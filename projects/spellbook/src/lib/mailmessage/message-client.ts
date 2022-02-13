
import {Observable, Subscriber} from 'rxjs';
import {MailMessage} from './mailmessage';
import {map} from 'rxjs/operators';
import {Spellbook} from '../core/spellbook';
import {Filter, ListResponse} from '../core/client';

export class MessageClient {

	constructor(private spellbook: Spellbook) {}

	// MailMessage

	public getMailMessageList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<ListResponse<MailMessage>> {
		const query = this.spellbook.api.createQueryList(page, results, filters, orderField, order);
		return this.spellbook.api.get<ListResponse<MailMessage>>(`/mailmessage${query}`).pipe(
			map((res: any) => {
				const items = new Array<MailMessage>();
				const response = res.items;
				for (const p of response) {
					items.push(p);
				}
				return new ListResponse(items, res.more);
			})
		);
	}

	public getAllMailMessage(): Observable<Array<MailMessage>> {
		let observer: Subscriber<Array<MailMessage>>;
		return new Observable<Array<MailMessage>>((ob) => {
			observer = ob;
			this.recursiveMailMessage([], observer, 0, 100);
		});
	}

	private recursiveMailMessage(all: Array<MailMessage>, observer: Subscriber<Array<MailMessage>>, page: number, results: number): void {
		this.getMailMessageList(page, results).subscribe((list: ListResponse<MailMessage>) => {
				all = all.concat(list.items);
				if (list.more) {
					this.recursiveMailMessage(all, observer, page + 1, results);
				} else {
					observer.next(all);
					observer.complete();
				}
			}
		);
	}

	public getMailMessage(id: number): Observable<MailMessage> {
		return this.spellbook.api.get<MailMessage>(`/mailmessage/${id}`).pipe(
			map((res: any) => {
				return new MailMessage(res);
			})
		);
	}

	public deleteMailMessage(mailmessage: MailMessage): Observable<void> {
		return this.spellbook.api.delete<void>(`/mailmessage/${mailmessage.id}`).pipe(
			map(() => {
				return;
			})
		);
	}

	public downloadCsv(): Observable<boolean> {
		// const filters: Filter[] = [];
		// let query = this.createQueryList(0, 999, filters);
		// query += `&property=Recipient`;
		const query = `?property=Recipient`;
		return this.spellbook.api.getFile<any>(`/mailmessage${query}`).pipe(
			map((res: any) => {
				console.log('res', res);
				const downloadURL = window.URL.createObjectURL(res);
				const link = document.createElement('a');
				link.href = downloadURL;
				link.download = 'mailmessage.csv';
				link.click();

				return true;
			})
		);
	}

	public getMailMessageProperties(property: string): Observable<string[]> {
		const query = `?property=${property}`;
		return this.spellbook.api.get<string[]>(`/mailmessage${query}`).pipe(
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
