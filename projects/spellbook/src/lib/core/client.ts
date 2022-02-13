import {HttpClient, HttpErrorResponse, HttpHeaders} from '@angular/common/http';
import {BehaviorSubject, Observable, throwError} from 'rxjs';
import {catchError, map} from 'rxjs/operators';
import {User} from './user';
import {FileAttachment} from '../media/multimedia/file';

export class Client {
	private static readonly SERVER_TOKEN_HEADER = 'X-Authentication';
	public static readonly LOCAL_TOKEN_KEY = 'bearer';

	public static readonly orderAscKey = 'asc';
	public static readonly orderDescKey = 'desc';

	private errorSubject: BehaviorSubject<HttpErrorResponse>;
	private readonly errorChannel: Observable<HttpErrorResponse>;

	constructor(private readonly apiUrl: string, private http: HttpClient) {
		this.errorSubject = new BehaviorSubject<HttpErrorResponse>(null);
		this.errorChannel = this.errorSubject.asObservable();
	}

	get networkError(): Observable<HttpErrorResponse> {
		return this.errorChannel;
	}

	private get authorizationHeaders(): HttpHeaders {
		let headers: HttpHeaders = new HttpHeaders();
		const token: string = localStorage.getItem(Client.LOCAL_TOKEN_KEY);
		if (token) {
			headers = headers.append(Client.SERVER_TOKEN_HEADER, token);
		}
		return headers;
	}

	private handleError(error: HttpErrorResponse | any) {
		console.error('[Network Error]:', error);
		if (error instanceof HttpErrorResponse) {
			const e = error as HttpErrorResponse;
			this.errorSubject.next(e);
			return throwError(e);
		}
		return throwError(error);
	}

	private handleRequest<T>(req: Observable<T>): Observable<T> {
		return req.pipe(catchError(
			(error: HttpErrorResponse) => {
				return this.handleError(error);
			}
		));
	}

	public get<T>(endpoint: string): Observable<T> {
		const headers = this.authorizationHeaders;
		const req = this.http.get<T>(
			`${this.apiUrl}${endpoint}`,
			{
				headers
			}
		);
		return this.handleRequest<T>(req);
	}

	public head(endpoint: string): Observable<any> {
		const headers = this.authorizationHeaders;
		const req = this.http.head(
			`${endpoint}`,
			{
				headers,
				observe: 'response'
			}
		);
		return this.handleRequest(req);
	}

	public getFile<T>(endpoint: string): Observable<T> {
		// const req = this.get(endpoint);
		const headers = this.authorizationHeaders;
		const responseType = 'blob' as 'json';
		const req = this.http.get<T>(
			`${this.apiUrl}${endpoint}`,
			{
				headers,
				responseType
			}
		);
		return this.handleRequest<T>(req);
	}

	public getList<T>(endpoint: string): Observable<ListResponse<T>> {
		return this.get<T>(endpoint).pipe(
			map((res: any) => {
				const listResponse = new ListResponse<T>(new Array<T>(), res.more);
				const items = res.items;
				for (const i of items) {
					items.push(i);
				}
				listResponse.items = items;
				return listResponse;
			})
		);
	}

	public post<T>(endpoint: string, body?: any): Observable<T> {
		let headers = this.authorizationHeaders;
		headers = headers.append('Content-Type', 'application/json');
		const req = this.http.post<T>(
			`${this.apiUrl}${endpoint}`,
			body,
			{
				headers
			}
		);
		return this.handleRequest(req);
	}


	public put<T>(endpoint: string, body?: any): Observable<T> {
		let headers = this.authorizationHeaders;
		headers = headers.append('Content-Type', 'application/json');
		const req = this.http.put<T>(
			`${this.apiUrl}${endpoint}`,
			body,
			{
				headers
			}
		);
		return this.handleRequest(req);
	}

	public delete<T>(endpoint: string): Observable<T> {
		let headers = this.authorizationHeaders;
		headers = headers.append('Content-Type', 'application/json');
		const req = this.http.delete<T>(
			`${this.apiUrl}${endpoint}`,
			{
				headers
			}
		);
		return this.handleRequest(req);
	}

	public postFile(type: string, namespace: string, name: string, body: Blob): Observable<FileAttachment> {
		const formData: FormData = new FormData();
		formData.append('file', body);
		formData.append('type', type);
		formData.append('namespace', namespace);
		formData.append('name', name);
		const headers = this.authorizationHeaders;
		// headers = headers.append('Content-Type', 'multipart/form-data');
		return this.http.post(`${this.apiUrl}/file`, formData, {headers}).pipe(
			map((res: any) => {
				return new FileAttachment(res);
			})
		);
	}

	public getUrlFiles(page?: number, results?: number): Observable<ListResponse<FileAttachment>> {
		const query = `?page=${page ? page : 0}&results=${results ? results : 100}`;
		return this.get<string[]>(`/file${query}`).pipe(
			map((res: any) => {
				const items = new Array<FileAttachment>();
				const iteresponses = res.items;
				for (const p of iteresponses) {
					items.push(p);
				}
				return new ListResponse(items, res.more);
			})
		);
	}

	public createQueryList(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): string {
		let query = `?page=${page ? page : 0}
		&results=${results ? results : 50}`;
		if (order && orderField) {
			query += `&order=`;
			if (order === Client.orderDescKey) {
				query += `-`;
			}
			query += `${orderField}`;
		}
		let i = 0;
		if (filters) {
			filters.forEach((f: Filter) => {
				if (i === 0) {
					query = query + '&filter=';
					query = query + `${f.field}=${f.value}`;
				} else {
					query = query + `^${f.field}=${f.value}`;
				}
				i++;
			});
		}
		return query;
	}

	/**
	 * User handling methods
	 *
	 */

	public createUser(user: User): Observable<User> {
		return this.post<User>('/users', user).pipe(
			map((res: any) => {
				return new User(res);
			})
		);
	}

	public updateUser(user: User): Observable<User> {
		return this.put<User>(`/users/${user.username}`, user).pipe(
			map((res: any) => {
				return new User(res);
			})
		);
	}

	public getUser(username: string): Observable<User> {
		return this.get<User>(`/users/${username}`).pipe(
			map((res: any) => {
				return new User(res);
			})
		);
	}

	public getUsers(page?: number, results?: number, filters?: Filter[], orderField?: string, order?: string): Observable<any> {
		const query = this.createQueryList(page, results, filters, orderField, order);
		return this.get<User[]>(`/users${query}`).pipe(
			map((res: any) => {
				const users = new Array<User>();
				const items = res.items;
				for (const u of items) {
					users.push(u);
				}
				return {items: users, more: res.more};
			})
		);
	}

	public createToken(username: string, password: string): Observable<string> {
		const params = {
			username,
			password
		};

		return this.post<string>('/tokens', params).pipe(
			map((res: string) => {
				return res;
			})
		);
	}

	public deleteToken(username: string): Observable<void> {
		return this.delete<void>(`/tokens/${username}`).pipe(
			map(() => {
				return;
			})
		);
	}

	public me(): Observable<User> {
		return this.get('/me').pipe(
			map((res: any) => {
				return new User(res);
			})
		);
	}

	/**
	 * LANGUAGES
	 */

	public getLanguages(): Observable<string[]> {
		return this.get<string[]>('/languages').pipe(
			map((res: any) => {
				const languages = new Array<string>();
				const items = res.items;
				if (items) {
					for (const i of items) {
						languages.push(i);
					}
				}
				return languages;
			})
		);
	}
}

export class ListResponse<Item> {
	public items: Array<Item>;
	public more: boolean;

	constructor(items: Item[], more: boolean) {
		this.items = items;
		this.more = more;
	}
}

export class Filter {
	public field: string;
	public value: string;

	constructor(field: string, value: string) {
		this.field = field;
		this.value = value;
	}
}
