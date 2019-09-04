package content

import (
	"context"
	"decodica.com/spellbook"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	cloudtasks "google.golang.org/api/cloudtasks/v2beta3"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type Task struct {
	CloudTask     *cloudtasks.Task `model:"-"`
	Name          string           `model:"-"`
	Url           string           `model:"-"`
	ScheduleTime  string           `model:"-"`
	ResponseTime  string           `model:"-"`
	Message       string           `model:"-"`
	DispatchCount int64            `model:"-"` // tentativi con errore
	ResponseCount int64            `model:"-"` // esecuzioni

	Method string `model:"-"`
}

func (task *Task) UnmarshalJSON(data []byte) error {

	alias := struct {
		Name          string `json:"name"`
		Url           string `json:"url"`
		ScheduleTime  string `json:"scheduleTime"`
		ResponseTime  string `json:"responseTime"`
		Message       string `json:"message"`
		DispatchCount int64  `json:"dispatchCount"`
		ResponseCount int64  `json:"responseCount"`

		Method string `json:"method"`
	}{}

	err := json.Unmarshal(data, &alias)
	if err != nil {
		return err
	}

	task.Name = alias.Name
	task.Url = alias.Url
	task.ScheduleTime = alias.ScheduleTime
	task.ResponseTime = alias.ResponseTime
	task.Message = alias.Message
	task.DispatchCount = alias.DispatchCount
	task.ResponseCount = alias.ResponseCount
	task.Method = alias.Method
	return nil
}

func (task *Task) MarshalJSON() ([]byte, error) {
	alias := struct {
		Name          string `json:"name"`
		Url           string `json:"url"`
		ScheduleTime  string `json:"scheduleTime"`
		ResponseTime  string `json:"responseTime"`
		Message       string `json:"message"`
		DispatchCount int64  `json:"dispatchCount"`
		ResponseCount int64  `json:"responseCount"`

		Method string `json:"method"`
	}{}

	alias.Name = task.CloudTask.Name
	alias.Url = task.CloudTask.AppEngineHttpRequest.RelativeUri
	if task.CloudTask.LastAttempt != nil {
		alias.ScheduleTime = task.CloudTask.LastAttempt.ScheduleTime
		alias.ResponseTime = task.CloudTask.LastAttempt.ResponseTime
		alias.Message = task.CloudTask.LastAttempt.ResponseStatus.Message
	}

	alias.DispatchCount = task.CloudTask.DispatchCount
	alias.ResponseCount = task.CloudTask.ResponseCount
	alias.Method = task.CloudTask.AppEngineHttpRequest.HttpMethod
	return json.Marshal(&alias)
}

/**
* Resource representation
 */

func (task *Task) Id() string {
	return task.CloudTask.Name
}

func (task *Task) FromRepresentation(rtype spellbook.RepresentationType, data []byte) error {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Unmarshal(data, task)
	}
	return spellbook.NewUnsupportedError()
}

func (task *Task) ToRepresentation(rtype spellbook.RepresentationType) ([]byte, error) {
	switch rtype {
	case spellbook.RepresentationTypeJSON:
		return json.Marshal(task)
	}
	return nil, spellbook.NewUnsupportedError()
}

func NewTaskController(projectId string, locationId string, queueId string) *spellbook.RestController {
	man := taskManager{}
	initTaskManager(&man, projectId, locationId, queueId)
	return spellbook.NewRestController(spellbook.BaseRestHandler{Manager: man})
}

func NewTaskControllerWithKey(key string, projectId string, locationId string, queueId string) *spellbook.RestController {
	man := taskManager{}
	initTaskManager(&man, projectId, locationId, queueId)
	handler := spellbook.BaseRestHandler{Manager: man}
	c := spellbook.NewRestController(handler)
	c.Key = key
	return c
}

func initTaskManager(man *taskManager, projectId string, locationId string, queueId string) {
	man.projectid = projectId
	man.locationid = locationId
	man.queueid = queueId
	if appengine.IsDevAppServer() {
		// projects/mage-middleware/locations/europe-west1/queues/spellbook-queue
		man.projectid = "mage-middleware"
		man.locationid = "europe-west1"
		man.queueid = "spellbook-queue"
	}
}

/*
* Task manager
 */

type taskManager struct {
	projectid  string
	locationid string
	queueid    string
}

func (manager taskManager) NewResource(ctx context.Context) (spellbook.Resource, error) {
	return &Task{}, nil
}

func (manager taskManager) FromId(ctx context.Context, id string) (spellbook.Resource, error) {
	if current := spellbook.IdentityFromContext(ctx); current == nil || !current.HasPermission(spellbook.PermissionReadAction) {
		return nil, spellbook.NewPermissionError(spellbook.PermissionName(spellbook.PermissionReadAction))
	}

	cloudtasksService, err := manager.GetCloudTasksService(ctx)
	if err != nil {
		return nil, spellbook.NewFieldError("cloudtasksService", err)
	}

	name := fmt.Sprintf("projects/%s/locations/%s/queues/%s/tasks/%s", manager.projectid, manager.locationid, manager.queueid, id)

	cloudTask, err := cloudtasksService.Projects.Locations.Queues.Tasks.Get(name).Context(ctx).Do()
	if err != nil {
		return nil, spellbook.NewFieldError("task get", err)
	}

	var t Task
	t.CloudTask = cloudTask
	return &t, nil
	//return nil, spellbook.NewUnsupportedError()
}

func (manager taskManager) ListOf(ctx context.Context, opts spellbook.ListOptions) ([]spellbook.Resource, error) {
	if current := spellbook.IdentityFromContext(ctx); current == nil || !current.HasPermission(spellbook.PermissionReadAction) {
		return nil, spellbook.NewPermissionError(spellbook.PermissionName(spellbook.PermissionReadAction))
	}

	cloudtasksService, err := manager.GetCloudTasksService(ctx)
	if err != nil {
		return nil, spellbook.NewFieldError("cloudtasksService", err)
	}

	// get queue
	myQueue, err := manager.GetQueue(ctx, cloudtasksService)
	if myQueue == nil {
		return nil, spellbook.NewFieldError("queue create", err)
	}

	// get tasks
	//var tasks []Task
	tasks := make([]Task, 0, 0)

	parent := fmt.Sprintf("projects/%s/locations/%s/queues/%s", manager.projectid, manager.locationid, manager.queueid)
	reqTasks := cloudtasksService.Projects.Locations.Queues.Tasks.List(parent)
	if err = reqTasks.Pages(ctx, func(page *cloudtasks.ListTasksResponse) error {
		for _, task := range page.Tasks {
			t := Task{}
			t.CloudTask = task
			tasks = append(tasks, t)
		}
		return nil
	}); err != nil {
		return nil, spellbook.NewFieldError("task list", err)
	}

	from := opts.Page * opts.Size
	if from > len(tasks) {
		return make([]spellbook.Resource, 0), nil
	}

	to := from + opts.Size
	if to > len(tasks) {
		to = len(tasks)
	}

	items := tasks[from:to]
	resources := make([]spellbook.Resource, len(items))

	for i := range items {
		task := Task(items[i])
		resources[i] = spellbook.Resource(&task)
	}

	return resources, nil
}

func (manager taskManager) ListOfProperties(ctx context.Context, opts spellbook.ListOptions) ([]string, error) {
	return nil, spellbook.NewUnsupportedError()
}

func (manager taskManager) Create(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	if current := spellbook.IdentityFromContext(ctx); current == nil || !current.HasPermission(spellbook.PermissionWriteAction) {
		return spellbook.NewPermissionError(spellbook.PermissionName(spellbook.PermissionWriteAction))
	}

	task := res.(*Task)
	task.CloudTask = &cloudtasks.Task{}
	task.CloudTask.AppEngineHttpRequest = &cloudtasks.AppEngineHttpRequest{}
	task.CloudTask.AppEngineHttpRequest.HttpMethod = task.Method
	task.CloudTask.AppEngineHttpRequest.RelativeUri = task.Url
	task.CloudTask.Name = fmt.Sprintf("projects/%s/locations/%s/queues/%s/tasks/%s", manager.projectid, manager.locationid, manager.queueid, task.Name)
	log.Infof(ctx, "task %+v", task.CloudTask)

	cloudtasksService, err := manager.GetCloudTasksService(ctx)
	if err != nil {
		return spellbook.NewFieldError("cloudtasksService", err)
	}

	myQueue, err := manager.GetQueue(ctx, cloudtasksService)
	if myQueue == nil {
		return spellbook.NewFieldError("queue not found", err)
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/queues/%s", manager.projectid, manager.locationid, manager.queueid)
	rb := &cloudtasks.CreateTaskRequest{
		Task: task.CloudTask,
	}
	cloudTask, err := cloudtasksService.Projects.Locations.Queues.Tasks.Create(parent, rb).Context(ctx).Do()
	if err != nil {
		return spellbook.NewFieldError("create task", err)
	}

	task.CloudTask = cloudTask

	return nil
}

func (manager taskManager) Update(ctx context.Context, res spellbook.Resource, bundle []byte) error {
	if current := spellbook.IdentityFromContext(ctx); current == nil || !current.HasPermission(spellbook.PermissionWriteAction) {
		return spellbook.NewPermissionError(spellbook.PermissionName(spellbook.PermissionWriteAction))
	}

	// RUN
	task := res.(*Task)
	task.CloudTask = &cloudtasks.Task{}
	task.CloudTask.AppEngineHttpRequest = &cloudtasks.AppEngineHttpRequest{}
	task.CloudTask.AppEngineHttpRequest.RelativeUri = task.Url
	task.CloudTask.Name = fmt.Sprintf("projects/%s/locations/%s/queues/%s/tasks/%s", manager.projectid, manager.locationid, manager.queueid, task.Name)
	log.Infof(ctx, "task %+v", task.CloudTask)

	cloudtasksService, err := manager.GetCloudTasksService(ctx)
	if err != nil {
		return spellbook.NewFieldError("cloudtasksService", err)
	}

	myQueue, err := manager.GetQueue(ctx, cloudtasksService)
	if myQueue == nil {
		return spellbook.NewFieldError("queue not found", err)
	}

	name := fmt.Sprintf("projects/%s/locations/%s/queues/%s/tasks/%s", manager.projectid, manager.locationid, manager.queueid, task.Name)
	rb := &cloudtasks.RunTaskRequest{}
	cloudTask, err := cloudtasksService.Projects.Locations.Queues.Tasks.Run(name, rb).Context(ctx).Do()
	if err != nil {
		return spellbook.NewFieldError("run task", err)
	}
	task.CloudTask = cloudTask

	return spellbook.NewUnsupportedError()
}

func (manager taskManager) Delete(ctx context.Context, res spellbook.Resource) error {
	return spellbook.NewUnsupportedError()
}

/**
* Utils
 */

func (manager taskManager) GetCloudTasksService(ctx context.Context) (*cloudtasks.Service, error) {
	c, err := google.DefaultClient(ctx, cloudtasks.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	cloudtasksService, err := cloudtasks.New(c)
	if err != nil {
		return nil, err
	}
	return cloudtasksService, nil
}

func (manager taskManager) GetQueue(ctx context.Context, cloudtasksService *cloudtasks.Service) (*cloudtasks.Queue, error) {
	var myQueue *cloudtasks.Queue
	parent := fmt.Sprintf("projects/%s/locations/%s", manager.projectid, manager.locationid)
	reqQueues := cloudtasksService.Projects.Locations.Queues.List(parent)
	if err := reqQueues.Pages(ctx, func(page *cloudtasks.ListQueuesResponse) error {
		for _, queue := range page.Queues {
			myQueue = queue
			break
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if myQueue == nil {
		return nil, errors.New("queue not found")

		// create queue
		/*parent := fmt.Sprintf("projects/%s/locations/%s", manager.projectid, manager.locationid)
		rb := &cloudtasks.Queue{
		}
		_, err := cloudtasksService.Projects.Locations.Queues.Create(parent, rb).Context(ctx).Do()
		if err != nil {
			return nil, errors.New(fmt.Sprintf("queue not created %s", err.Error()))
		}
		return manager.GetQueue(ctx, cloudtasksService)*/
	}
	return myQueue, nil
}
