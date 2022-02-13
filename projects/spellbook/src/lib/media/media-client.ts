
import {Attachment} from './multimedia/attachment';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {Spellbook} from '../core/spellbook';
import {Filter, ListResponse} from '../core/client';

export class MediaClient {

	constructor(private spellbook: Spellbook) {}

	/**
	 * MULTIMEDIA
	 */

	public createAttachment(attachment: Attachment): Observable<Attachment> {
		return this.spellbook.api.post<Attachment>('/attachment', attachment).pipe(
			map((res: any) => {
				return new Attachment(res);
			})
		);
	}

	public updateAttachment(multimedia: Attachment): Observable<Attachment> {
		return this.spellbook.api.put<Attachment>(`/attachment/${multimedia.id}`, multimedia).pipe(
			map((res: any) => {
				return new Attachment(res);
			})
		);
	}

	public getAttachment(id: string): Observable<Attachment> {
		return this.spellbook.api.get<Attachment>(`/attachment/${id}`).pipe(
			map((res: any) => {
				return new Attachment(res);
			})
		);
	}

	public getAttachmentList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<ListResponse<Attachment>> {
		const query = this.spellbook.api.createQueryList(page, results, filters, orderField, order);
		return this.spellbook.api.get<ListResponse<Attachment>>(`/attachment${query}`).pipe(
			map((res: any) => {
				const items = new Array<Attachment>();
				const response = res.items;
				for (const p of response) {
					items.push(p);
				}
				return new ListResponse(items, res.more);
			})
		);
	}

	public getAttachmentPropertyList(property: string): Observable<string[]> {
		const query = `?property=${property}`;
		return this.spellbook.api.get<string[]>(`/attachment${query}`).pipe(
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

	public deleteAttachment(attachment: Attachment): Observable<void> {
		return this.spellbook.api.delete<void>(`/attachment/${attachment.id}`).pipe(
			map(() => {
				return;
			})
		);
	}
}
