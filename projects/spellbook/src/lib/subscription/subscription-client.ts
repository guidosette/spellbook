import {Observable, Subscriber} from 'rxjs';
import {Subscription} from './subscription';
import {map} from 'rxjs/operators';
import {Spellbook} from '../core/spellbook';
import {Filter, ListResponse} from '../core/client';

export class SubscriptionClient {

	constructor(private spellbook: Spellbook) {
	}

	// Subscription

	public getSubscriptionList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<ListResponse<Subscription>> {
		const query = this.spellbook.api.createQueryList(page, results, filters, orderField, order);
		return this.spellbook.api.get<ListResponse<Subscription>>(`/subscription${query}`).pipe(
			map((res: any) => {
				const items = new Array<Subscription>();
				const response = res.items;
				for (const p of response) {
					items.push(p);
				}
				return new ListResponse(items, res.more);
			})
		);
	}

	public getAllSubscription(): Observable<Array<Subscription>> {
		let observer: Subscriber<Array<Subscription>>;
		return new Observable<Array<Subscription>>((ob) => {
			observer = ob;
			this.recursiveSubscription([], observer, 0, 100);
		});
	}

	private recursiveSubscription(all: Array<Subscription>, observer: Subscriber<Array<Subscription>>, page: number, results: number): void {
		this.getSubscriptionList(page, results).subscribe((list: ListResponse<Subscription>) => {
				all = all.concat(list.items);
				if (list.more) {
					this.recursiveSubscription(all, observer, page + 1, results);
				} else {
					observer.next(all);
					observer.complete();
				}
			}
		);
	}

	public getSubscription(key: string): Observable<Subscription> {
		return this.spellbook.api.get<Subscription>(`/subscription/${key}`).pipe(
			map((res: any) => {
				return new Subscription(res);
			})
		);
	}

	public deleteSubscription(subscription: Subscription): Observable<void> {
		return this.spellbook.api.delete<void>(`/subscription/${subscription.key}`).pipe(
			map(() => {
				return;
			})
		);
	}

	public downloadCsv(): Observable<boolean> {
		// const filters: Filter[] = [];
		// let query = this.spellbook.api.createQueryList(0, 999, filters);

		// return this.spellbook.api.get<ListResponse<Subscription>>(`/subscription${query}`).pipe(
		// 	map((res: any) => {
		// 		const items = new Array<Subscription>();
		// 		const response = res.items;
		// 		for (const p of response) {
		// 			items.push(p);
		// 		}
		// 		return new ListResponse(items, res.more);
		// 	})
		// );

		// todo
		const query = `?property=Email`;
		return this.spellbook.api.getFile<any>(`/subscription${query}`).pipe(
			map((res: any) => {
				console.log('res', res);
				const downloadURL = window.URL.createObjectURL(res);
				const link = document.createElement('a');
				link.href = downloadURL;
				link.download = 'subscription.csv';
				link.click();

				return true;
			})
		);
	}

	public getSubscriptionProperties(property: string): Observable<string[]> {
		const query = `?property=${property}`;
		return this.spellbook.api.get<string[]>(`/subscription${query}`).pipe(
			map((res: any) => {
				const result = new Array<string>();
				const items = res.items;
				if (items) {
					for (const p of items) {
						result.push(p);
					}
				}
				return result;
			})
		);
	}
}
