import {Page} from './page';
import {Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {Spellbook} from '../core/spellbook';
import {Filter, ListResponse} from '../core/client';

export class NavigationClient {

	constructor(private spellbook: Spellbook) {}

	/**
	 * NAVIGATION
	 */
	public createPage(page: Page): Observable<Page> {
		return this.spellbook.api.post<Page>('/page', page).pipe(
			map((res: any) => {
				return new Page(res);
			})
		);
	}

	public updatePage(page: Page): Observable<Page> {
		return this.spellbook.api.put<Page>(`/page/${page.id}`, page).pipe(
			map((res: any) => {
				return new Page(res);
			})
		);
	}

	public deletePage(page: Page): Observable<void> {
		return this.spellbook.api.delete<void>(`/page/${page.id}`).pipe(
			map(() => {
				return;
			})
		);
	}

	public getPage(id: number): Observable<Page> {
		return this.spellbook.api.get<Page>(`/page/${id}`).pipe(
			map((res: any) => {
				return new Page(res);
			})
		);
	}

	public getPageList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<ListResponse<Page>> {
		const query = this.spellbook.api.createQueryList(page, results, filters, orderField, order);
		return this.spellbook.api.get<ListResponse<Page>>(`/page${query}`).pipe(
			map((res: any) => {
				const items = new Array<Page>();
				const response = res.items;
				for (const s of response) {
					items.push(s);
				}
				return new ListResponse(items, res.more);
			})
		);
	}


	public getStaticPageCodeList(): Observable<string[]> {
		const query = ``;
		return this.spellbook.api.get<string[]>(`/staticpage${query}`).pipe(
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
