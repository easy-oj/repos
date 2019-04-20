package submit

import (
	"context"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/proto/queue"
	"github.com/easy-oj/repos/common/caller"
	"github.com/easy-oj/repos/common/database"
)

func Submit(uuid string, content map[string]string) (int32, error) {
	var uid, pid, rid, lid, sid int32
	row := database.DB.QueryRow("SELECT uid, pid, id, lid FROM tb_repo WHERE uuid = ?", uuid)
	if err := row.Scan(&uid, &pid, &rid, &lid); err != nil {
		return sid, err
	}
	res, err := database.DB.Exec(
		"INSERT INTO tb_submission (uid, pid, rid, lid, executions) VALUES (?, ?, ?, ?, ?)", uid, pid, rid, lid, "[]")
	if err != nil {
		return sid, err
	}
	if id, err := res.LastInsertId(); err != nil {
		return sid, err
	} else {
		sid = int32(id)
	}
	_, err = database.DB.Exec("UPDATE tb_problem SET submitted_count = submitted_count + 1 WHERE id = ?", pid)
	if err != nil {
		logs.Error("[Submit] sid = %d, update submitted count error: %s", sid, err.Error())
	}
	req := queue.NewPutMessageReq()
	req.Message = &queue.Message{
		Uid:     uid,
		Pid:     pid,
		Lid:     lid,
		Sid:     sid,
		Content: content,
	}
	_, err = caller.QueueClient.PutMessage(context.Background(), req)
	return sid, err
}
