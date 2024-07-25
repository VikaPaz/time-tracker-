package task

import (
	"database/sql"
	"errors"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type TaskRepository struct {
	conn *sql.DB
	log  *logrus.Logger
}

func NewRepository(conn *sql.DB, logger *logrus.Logger) *TaskRepository {
	return &TaskRepository{
		conn: conn,
		log:  logger,
	}
}

func (r *TaskRepository) Create(task models.Task) (models.Task, error) {
	r.log.Debugf("Executing insert task: %+v", task)
	row := r.conn.QueryRow("INSERT INTO tasks (task, user_id) VALUES "+
		"($1, $2) RETURNING id", task.Task, task.UserID)
	if err := row.Err(); err != nil {
		return models.Task{}, models.ErrCreateTaskResponse
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return models.Task{}, models.ErrCreateTaskResponse
	}
	task.ID = id
	return task, nil
}

func (r *TaskRepository) Get(request models.LaborTimeRequest) (models.LaborTimeResponse, error) {
	r.log.Debugf("Executing query")
	rows, err := r.conn.Query(`select t.id,
       t.task,
       extract(epoch from (sum(l.stop - l.start)))::int as delta
from tasks t
         join public.labor_time l on t.id = l.task_id
where t.user_id = $3
  and l.start between $2
    and $1
  and l.stop <= $1
group by t.id
order by delta`,
		request.EndTime, request.StartTime, request.UserID)
	if err != nil {
		return models.LaborTimeResponse{}, errors.Join(models.ErrGetTaskResponse, err)
	}

	resp := models.LaborTimeResponse{UserID: *request.UserID}
	for rows.Next() {
		task := models.TaskInfo{}

		var duration time.Duration
		err = rows.Scan(&task.ID, &task.Task, &duration)
		if err != nil {
			return models.LaborTimeResponse{}, errors.Join(models.ErrGetTaskResponse, err)
		}
		task.LaborTime = &duration
		resp.Tasks = append(resp.Tasks, task)
	}
	return resp, nil
}

func (r *TaskRepository) Start(taskID uuid.UUID) error {
	r.log.Debugf("Executing query")
	_, err := r.conn.Exec("INSERT INTO labor_time (task_id) VALUES ($1)", taskID)
	if err != nil {
		return models.ErrStartTimer
	}
	return nil
}

func (r *TaskRepository) Stop(taskID uuid.UUID) error {
	r.log.Debugf("Executing query")
	_, err := r.conn.Exec("update labor_time set stop = now() WHERE task_id = $1 and stop is null", taskID)
	if err != nil {
		return models.ErrStopTimer
	}
	return nil
}

func (r *TaskRepository) IsStarted(taskID uuid.UUID) (bool, error) {
	r.log.Debugf("Executing query")
	row := r.conn.QueryRow("SELECT id FROM labor_time WHERE task_id = $1 and stop is null", taskID)
	if err := row.Err(); err != nil {
		return false, models.ErrCheckTimerStatus
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return false, models.ErrCheckTimerStatus
	}

	if id == uuid.Nil {
		return false, nil
	}
	return true, models.ErrTimerStarted
}
