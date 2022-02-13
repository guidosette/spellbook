import {Spellbook} from '../core/spellbook';
import {Observable} from 'rxjs';
import {Content} from './content';
import {map} from 'rxjs/operators';
import {Filter, ListResponse} from '../core/client';

export class ContentClient {

	constructor(private spellbook: Spellbook) {
	}

	/**
	 * CONTENTS
	 */
	public createContent(content: Content): Observable<Content> {
		return this.spellbook.api.post<Content>('/content', content).pipe(
			map((res: any) => {
				return new Content(res, this.spellbook);
			})
		);
	}

	public updateContent(content: Content): Observable<Content> {
		return this.spellbook.api.put<Content>(`/content/${content.id}`, content).pipe(
			map((res: any) => {
				return new Content(res, this.spellbook);
			})
		);
	}

	public deleteContent(content: Content): Observable<void> {
		return this.spellbook.api.delete<void>(`/content/${content.id}`).pipe(
			map(() => {
				return;
			})
		);
	}

	public getContent(id: string | number): Observable<Content> {
		return this.spellbook.api.get<Content>(`/content/${id}`).pipe(
			map((res: any) => {
				return new Content(res, this.spellbook);
			})
		);
	}

	public getContentList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<ListResponse<Content>> {
		const query = this.spellbook.api.createQueryList(page, results, filters, orderField, order);
		return this.spellbook.api.get<ListResponse<Content>>(`/content${query}`).pipe(
			map((res: any) => {
				const items = new Array<Content>();
				const iteresponses = res.items;
				for (const p of iteresponses) {
					const c = new Content(p);
					items.push(c);
				}
				return new ListResponse(items, res.more);
			})
		);
	}

	public getContentProperties(property: string): Observable<string[]> {
		const query = `?property=${property}`;
		return this.spellbook.api.get<string[]>(`/content${query}`).pipe(
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

	public getContentSupportedLanguages(idTranslate: string): Observable<string[]> {
		const filters: Filter[] = [];
		filters.push(new Filter('IdTranslate', idTranslate));
		let query = this.spellbook.api.createQueryList(0, 50, filters);
		query += `&property=Locale`;
		return this.spellbook.api.get<string[]>(`/content${query}`).pipe(
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

	public getSpecialCodeList(type: string, category: string): Observable<string[]> {
		let query = `?`;
		query += `type=${type}`;
		query += `&category=${category}`;
		return this.spellbook.api.get<string[]>(`/specialcode${query}`).pipe(
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
