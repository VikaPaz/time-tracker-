package task

import (
	"database/sql"
	"github.com/VikaPaz/time_tracker/internal/models"
	"github.com/google/uuid"
	"time"
)

type TaskRepository struct {
	conn *sql.DB
}

func NewRepository(conn *sql.DB) *TaskRepository {
	return &TaskRepository{conn: conn}
}

func (r *TaskRepository) Create(task models.Task) (models.Task, error) {
	row := r.conn.QueryRow("INSERT INTO tasks (task, user_id) VALUES "+
		"($1, $2) RETURNING id", task.Task, task.UserID)
	if err := row.Err(); err != nil {
		return models.Task{}, err
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		return models.Task{}, err
	}

	task.ID = id
	return task, nil
}

func (r *TaskRepository) Get(request models.LaborTimeRequest) (models.LaborTimeResponse, error) {
	//res, err := r.conn.Query("select t.id, t.task, "+
	//	"(select sum(labor_time.stop - labor_time.start) "+
	//	"from labor_time where labor_time.task_id = t.id) as delta "+
	//	"from tasks t where t.user_id = $1 group by t.id", request.UserID)

	rows, err := r.conn.Query("select t.id, t.task, "+
		"extract(minutes from SUM(CASE WHEN l.stop IS NULL THEN now() - l.start ELSE l.stop - l.start END))::int AS delta "+
		"FROM tasks t LEFT JOIN labor_time l ON l.task_id = t.id AND l.start BETWEEN $1 AND $2 "+
		"AND (l.stop < $3 OR l.stop IS NULL) WHERE t.user_id = $4 "+
		"AND (CASE WHEN l.stop IS NULL THEN now() - l.start ELSE l.stop - l.start END) IS NOT NULL GROUP BY t.id;",
		request.Offset, request.Limit, request.Limit, request.UserID)
	if err != nil {
		return models.LaborTimeResponse{}, err
	}

	resp := models.LaborTimeResponse{UserID: *request.UserID, Total: 0}
	for rows.Next() {
		task := models.TaskInfo{}

		var duration time.Duration
		err = rows.Scan(&task.ID, &task.Task, &duration)
		duration *= time.Minute
		if err != nil {
			if err == sql.ErrNoRows {
				return models.LaborTimeResponse{}, nil
			}
			return models.LaborTimeResponse{}, err
		}
		task.LaborTime = &duration

		resp.Tasks = append(resp.Tasks, task)
		resp.Total += 1
	}

	return resp, nil
}

func (r *TaskRepository) Start(taskID uuid.UUID) error {
	_, err := r.conn.Exec("INSERT INTO labor_time (task_id) VALUES ($1)", taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) Stop(taskID uuid.UUID) error {
	//stop := time.Now()

	_, err := r.conn.Exec("update labor_time set stop = now() WHERE task_id = $1 and stop is null", taskID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) IsStarted(taskID uuid.UUID) (bool, error) {
	row := r.conn.QueryRow("SELECT id FROM labor_time WHERE task_id = $1 and stop is null", taskID)
	if err := row.Err(); err != nil {
		return false, err
	}

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	if id == uuid.Nil {
		return false, nil
	}
	return true, nil
}
