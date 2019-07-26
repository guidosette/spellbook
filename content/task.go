package content

import (
	"appengine"
	"context"
	"distudio.com/page"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	cloudtasks "google.golang.org/api/cloudtasks/v2beta3"
	"google.golang.org/appengine/log"
)

type Task struct {
	CloudTask *cloudtasks.Task `model:"-"`
}

func (task *Task) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name string `json:"name"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	task.CloudTask.Name = alias.Name

	return nil
}

func (task *Task) MarshalJSON() ([]byte, error) {
	alias := struct {
		Name string `json:"name"`
	}{task.CloudTask.Name}

	return json.Marshal(&alias)
}

/**
* Resource representation
 */

func (task *Task) Id() string {
	return task.CloudTask.Name
}

func (task *Task) FromRepresentation(rtype page.RepresentationType, data []byte) error {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Unmarshal(data, task)
	}
	return page.NewUnsupportedError()
}

func (task *Task) ToRepresentation(rtype page.RepresentationType) ([]byte, error) {
	switch rtype {
	case page.RepresentationTypeJSON:
		return json.Marshal(task)
	}
	return nil, page.NewUnsupportedError()
}

func NewTaskController(project_id string, location_id string, queue_id string) *page.RestController {
	man := taskManager{}
	man.projectid = project_id
	man.locationid = location_id
	man.queueid = queue_id
	if appengine.IsDevAppServer() {
		man.projectid = "mage-middleware"
		man.locationid = "europe-west1"
		man.queueid = "main-queue"
	}
	return page.NewRestController(page.BaseRestHandler{Manager: man})
}

/*
* Task manager
 */

type taskManager struct {
	projectid  string
	locationid string
	queueid    string
}

func (manager taskManager) NewResource(ctx context.Context) (page.Resource, error) {
	return nil, page.NewUnsupportedError()
}

func (manager taskManager) FromId(ctx context.Context, id string) (page.Resource, error) {
	c, err := google.DefaultClient(ctx, cloudtasks.CloudPlatformScope)
	if err != nil {
		return nil, page.NewFieldError("DefaultClient", err)
	}

	cloudtasksService, err := cloudtasks.New(c)
	if err != nil {
		return nil, page.NewFieldError("cloudtasks.New", err)
	}

	// projects/mage-middleware/locations/europe-west1/queues/my-queue
	parent := fmt.Sprintf("projects/%s/locations/%s/queues/%s", manager.projectid, manager.locationid, manager.queueid)

	req := cloudtasksService.Projects.Locations.Queues.Tasks.List(parent)
	if err := req.Pages(ctx, func(page *cloudtasks.ListTasksResponse) error {
		for _, task := range page.Tasks {
			// TODO: Change code below to process each `task` resource:
			fmt.Printf("%#v\n", task)
		}
		return nil
	}); err != nil {
		return nil, page.NewFieldError("task list", err)
	}

	return nil, page.NewUnsupportedError()
}

func (manager taskManager) ListOf(ctx context.Context, opts page.ListOptions) ([]page.Resource, error) {

	c, err := google.DefaultClient(ctx, cloudtasks.CloudPlatformScope)
	if err != nil {
		return nil, page.NewFieldError("DefaultClient", err)
	}

	cloudtasksService, err := cloudtasks.New(c)
	if err != nil {
		return nil, page.NewFieldError("cloudtasks.New", err)
	}

	// get queue
	var myQueue *cloudtasks.Queue
	// projects/mage-middleware/locations/europe-west1/queues/my-queue
	parent := fmt.Sprintf("projects/%s/locations/%s", manager.projectid, manager.locationid)

	reqQueues := cloudtasksService.Projects.Locations.Queues.List(parent)
	if err := reqQueues.Pages(ctx, func(page *cloudtasks.ListQueuesResponse) error {
		for _, queue := range page.Queues {
			myQueue = queue
			break
		}
		return nil
	}); err != nil {
		return nil, page.NewFieldError("queue list", err)
	}

	if myQueue == nil {
		return nil, page.NewFieldError("queue list", errors.New("queue not found"))
	}

	// get tasks
	var tasks []Task
	parent = fmt.Sprintf("projects/%s/locations/%s/queues/%s", manager.projectid, manager.locationid, manager.queueid)
	reqTasks := cloudtasksService.Projects.Locations.Queues.Tasks.List(parent)
	if err = reqTasks.Pages(ctx, func(page *cloudtasks.ListTasksResponse) error {
		tasks := make([]Task, len(page.Tasks))
		log.Debugf(ctx, "tasks: %+v", page.Tasks)
		for _, task := range page.Tasks {
			var t Task
			t.CloudTask = task
			tasks = append(tasks, t)
		}
		return nil
	}); err != nil {
		return nil, page.NewFieldError("task list", err)
	}

	from := opts.Page * opts.Size
	if from > len(tasks) {
		return make([]page.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(tasks) {
		to = len(tasks)
	}

	items := tasks[from:to]
	resources := make([]page.Resource, len(items))

	for i := range items {
		task := Task(items[i])
		resources[i] = page.Resource(&task)
	}

	return resources, nil
}

func (manager taskManager) ListOfProperties(ctx context.Context, opts page.ListOptions) ([]string, error) {
	return nil, page.NewUnsupportedError()
}

func (manager taskManager) Create(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager taskManager) Update(ctx context.Context, res page.Resource, bundle []byte) error {
	return page.NewUnsupportedError()
}

func (manager taskManager) Delete(ctx context.Context, res page.Resource) error {
	return page.NewUnsupportedError()
}
