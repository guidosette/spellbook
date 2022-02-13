import {Inject, Injectable, Injector} from '@angular/core';
import {BehaviorSubject, forkJoin, Observable} from 'rxjs';
import {User} from './user';
import {HttpClient, HttpErrorResponse} from '@angular/common/http';
import {SpellbookConfig} from './spellbook-config';
import {Client} from './client';
import {Router} from '@angular/router';
import {Category} from './category';
import {map} from 'rxjs/operators';
import {AttachmentType, SupportedAttachment} from './supported-attachment';
import {FileAttachment} from '../media/multimedia/file';
import {MatSnackBar} from '@angular/material';
import {Action} from '../action/action';
import {TypedefService} from './typedef.service';
import * as uuid from 'uuid';


@Injectable({
	providedIn: 'root'
})
export class Spellbook {

	private readonly userSubject: BehaviorSubject<User>;
	private readonly userObservable: Observable<User>;
	private readonly client: Client;
	private readonly typeDefService: TypedefService;
	public router: Router;
	private supportedCats: Category[] = [];
	private supportedActs: Action[] = [];

	constructor(@Inject('config') private config: SpellbookConfig, private injector: Injector, http: HttpClient, private snackBar: MatSnackBar) {
		this.userSubject = new BehaviorSubject<User>(null);
		this.userObservable = this.userSubject.asObservable();
		this.client = new Client(config.apiUrl, http);
		this.typeDefService = injector.get(TypedefService);
		this.client.networkError.subscribe(
			(err: HttpErrorResponse) => {
				this.handleNetworkError(err);
			}
		);
		for (const def of config.typeDefinitions()) {
			this.typeDefService.addTypeDefinition(def);
		}
	}

	public uniqueId(): string {
		return uuid.v4();
	}

	public get superUserUrl(): string {
		return this.config.superUserRedirectUrl;
	}

	private handleNetworkError(err: HttpErrorResponse) {
		if (!err) {
			return;
		}

		switch (err.status) {
			case 401:
				// got an unauthorize, leave the area!
				console.error('Unauthorized - log out');
				this.user = null;
				this.router.navigate(['/login']);
				break;
			case 404:
				console.error('Not Found', err);
				const s = err.url.split('/');
				if (s.length > 1) {
					const finalPath = s[s.length - 2];
					if (finalPath === 'content' || finalPath === 'users' || finalPath === 'place' || finalPath === 'media') {
						this.router.navigate(['/fourofour']);
					}
				}
				break;
			case 500:
				console.error('Internal Server error', err);
				this.snackBar.open('Error ' + err.statusText, 'ok', {});
				break;
		}
	}

	public get identity(): Observable<User> {
		return this.userObservable;
	}

	public set user(user: User) {
		this.userSubject.next(user);
	}

	public get user(): User {
		return this.userSubject.getValue();
	}

	public get definitions(): TypedefService {
		return this.typeDefService;
	}

	public set token(tkn: string) {
		this.user = null;
		if (tkn == null) {
			localStorage.removeItem(Client.LOCAL_TOKEN_KEY);
			return;
		}

		localStorage.setItem(Client.LOCAL_TOKEN_KEY, tkn);
	}

	public get token(): string {
		return localStorage.getItem(Client.LOCAL_TOKEN_KEY);
	}

	isLoggedIn(): boolean {
		return !!this.user;
	}

	public get api(): Client {
		return this.client;
	}

	public get supportedAttachments(): Array<SupportedAttachment> {
		return [
			{
				name: 'Gallery',
				value: AttachmentType.GALLERY,
				image: 'image',
				accept: '.png,.jpg,.jpeg,.svg,.bmp,.tiff,.gif',
				mime: 'image/png,image/jpeg,image/svg,image/bmp,image/tiff,image/gif'
			},
			{
				name: 'Attachment',
				value: AttachmentType.ATTACHMENT,
				image: 'attach_file',
				accept: '*/*',
				mime: null
			},
			{
				name: 'Video',
				value: AttachmentType.VIDEO,
				image: 'videocam',
				accept: '.webm,.mp4',
				mime: 'video/mp4,video/webm'
			},
		];
	}

	/**
	 * Return AttachmentType: "gallery", "attachments", "video"
	 */
	public getAttachmentTypeForFileAttachment(f: FileAttachment): string {
		let value = '';
		if (f.name) {
			value = f.name;
		} else {
			value = f.resourceUrl;
		}
		const s = value.split('.');
		let extension = '';
		if (s.length > 1) {
			extension = s.pop().toLowerCase();
		}
		const support = this.supportedAttachments.find((sa) => {
			return sa.accept.indexOf(extension) !== -1;
		});
		return support ? support.value : AttachmentType.ATTACHMENT;
	}

	public init(): Promise<any> {
		this.router = this.injector.get(Router);
		return new Promise((resolve, reject) => {
			const results = [];
			results.push(this.api.me());
			results.push(this.supportedCategories());
			forkJoin(results).subscribe((response: any) => {
					this.user = response[0];
					this.supportedCats = response[1];

					if (this.user.hasPermission(User.PERMISSION_READ_ACTION) || this.user.hasPermission(User.PERMISSION_WRITE_ACTION)) {
						this.supportedActs = new Array<Action>();
						this.supportedActions().subscribe((actions: Action[]) => {
								for (const a of actions) {
									this.supportedActs.push(new Action(a));
								}
								resolve();
							},
							(error: HttpErrorResponse) => {
								console.error('error init supportedActs', error);
								resolve();
							});
					} else {
						resolve();
					}
				},
				(error: HttpErrorResponse) => {
					console.error('error init', error);
					resolve();
				});
		});
	}

	/**
	 * Supported Categories
	 */
	private supportedCategories(): Observable<Category[]> {
		return this.api.get<Category[]>('/categories').pipe(
			map((res: any) => {
				const categories = new Array<Category>();
				const items = res.items;
				if (items) {
					for (const i of items) {
						categories.push(i);
					}
				}
				return categories;
			})
		);
	}

	public getSupportedCategories() {
		return this.supportedCats;
	}

	/**
	 * Supported Actions
	 */
	private supportedActions(): Observable<Action[]> {
		return this.api.get<Action[]>('/actions').pipe(
			map((res: any) => {
				const categories = new Array<Action>();
				const items = res.items;
				if (items) {
					for (const i of items) {
						categories.push(i);
					}
				}
				return categories;
			})
		);
	}

	public getSupportedActions() {
		return this.supportedActs;
	}
}
