import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {Filter, ListResponse} from '../core/client';
import {Action} from './action';
import {Spellbook} from '../core/spellbook';
import {Task} from './task';

export class ActionClient {

	constructor(private spellbook: Spellbook) {
	}

	/**
	 * ACTIONS
	 */

	public getActionList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<ListResponse<Action>> {
		return new Observable<ListResponse<Action>>((observer) => {
			observer.next(new ListResponse(this.spellbook.getSupportedActions(), false));
		});
	}

	/**
	 * TASKS
	 */

	public getTasks(): Observable<ListResponse<Task>> {
		const query = `?page=${0}&results=${100}`;
		return this.spellbook.api.get<string[]>(`/task${query}`).pipe(
			map((res: any) => {
				const items = new Array<Task>();
				const iteresponses = res.items;
				for (const p of iteresponses) {
					items.push(p);
				}
				return new ListResponse(items, res.more);
			})
		);
	}

	public createTask(task: Task): Observable<Task> {
		return this.spellbook.api.post<Task>('/task', task).pipe(
			map((res: any) => {
				return new Task(res);
			})
		);
	}

	public runTask(task: Task): Observable<Task> {
		return this.spellbook.api.put<Task>(`/task/${task.getSimpleName()}`, task).pipe(
			map((res: any) => {
				return new Task(res);
			})
		);
	}
}
