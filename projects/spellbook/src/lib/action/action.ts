import {Task} from './task';

export class Action {

	static readonly ACTION_TYPE_NORMAL: string = 'normal';
	static readonly ACTION_TYPE_UPLOAD: string = 'upload';

	public static readonly ACTION_STATUS_RUN = 'RUN';
	public static readonly ACTION_STATUS_STOP = 'STOP';

	name: string;
	endpoint: string;
	type: string;
	method: string;

	status: string;

	get task(): Task {
		return this._task;
	}
	set task(value: Task) {
		this._task = value;
		if (this._task && !this._task.responseTime) {
			this.status = Action.ACTION_STATUS_RUN;
		} else {
			this.status = Action.ACTION_STATUS_STOP;
		}
	}
	_task: Task;

	public isNew() {
		return !(this.endpoint && this.endpoint !== undefined && this.endpoint.length > 0);
	}

	public isRun() {
		return this.status === Action.ACTION_STATUS_RUN;
	}

	constructor(json?: any) {
		if (json) {
			this.endpoint = json.endpoint;
			this.name = json.name;
			this.type = json.type;
			this.method = json.method;

			this.status = Action.ACTION_STATUS_STOP;
		}
	}

	public toJSON(): any {
		return {
			endpoint: this.endpoint,
			name: this.name,
			status: this.status,
			type: this.type,
			method: this.method,
		};
	}

	public getTask() {
		if (this.isTaskCreated()) {
			return this.task;
		}
		const t: Task = new Task();
		t.name = this.name;
		t.url = this.endpoint;
		t.method = this.method;
		return t;
	}

	public isTaskCreated() {
		return this.task;
	}
}

