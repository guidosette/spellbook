
import {map} from 'rxjs/operators';
import {Observable} from 'rxjs';
import {Place} from './place';
import {Spellbook} from '../core/spellbook';
import {Filter, ListResponse} from '../core/client';

export class PlaceClient {

	constructor(private spellbook: Spellbook) {}

	// Place

	public createPlace(place: Place): Observable<Place> {
		return this.spellbook.api.post<Place>('/place', place).pipe(
			map((res: any) => {
				return new Place(res);
			})
		);
	}

	public updatePlace(place: Place): Observable<Place> {
		return this.spellbook.api.put<Place>(`/place/${place.id}`, place).pipe(
			map((res: any) => {
				return new Place(res);
			})
		);
	}

	public deletePlace(place: Place): Observable<void> {
		return this.spellbook.api.delete<void>(`/place/${place.id}`).pipe(
			map(() => {
				return;
			})
		);
	}

	public getPlace(id: number): Observable<Place> {
		return this.spellbook.api.get<Place>(`/place/${id}`).pipe(
			map((res: any) => {
				return new Place(res);
			})
		);
	}

	public getPlaceList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<ListResponse<Place>> {
		const query = this.spellbook.api.createQueryList(page, results, filters, orderField, order);
		return this.spellbook.api.get<ListResponse<Place>>(`/place${query}`).pipe(
			map((res: any) => {
				const items = new Array<Place>();
				const response = res.items;
				for (const p of response) {
					items.push(p);
				}
				return new ListResponse(items, res.more);
			})
		);
	}

	public getPlaceProperties(property: string): Observable<string[]> {
		const query = `?property=${property}`;
		return this.spellbook.api.get<string[]>(`/place${query}`).pipe(
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
