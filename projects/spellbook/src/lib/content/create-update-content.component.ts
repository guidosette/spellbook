import {Component, OnInit, ViewChild} from '@angular/core';
import {CkUploadAdapter} from '../core/ck-upload-adapter';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {ActivatedRoute} from '@angular/router';
import {Content} from './content';
import {Observable} from 'rxjs';
import {map, startWith} from 'rxjs/operators';
import {MatChipInputEvent, MatDialog, MatSnackBar} from '@angular/material';
import {COMMA, ENTER} from '@angular/cdk/keycodes';
import {HttpErrorResponse} from '@angular/common/http';
import * as ClassicEditor from '@ckeditor/ckeditor5-build-classic';

import {ContentClient} from './content-client';
import {AttachmentType, SupportedAttachment} from '../core/supported-attachment';
import {CreateAttachmentGroupComponent} from '../media/multimedia/create-attachment-group.component';
import {AttachmentGroup} from '../media/multimedia/attachment-group';
import {ListMediaDialogComponent} from '../media/list-media-dialog.component';
import {Attachment} from '../media/multimedia/attachment';
import {ResponseError} from '../core/response-error';
import {Spellbook} from '../core/spellbook';
import {Client, Filter, ListResponse} from '../core/client';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {UrlUtils} from '../core/url-utils';
import {ContentDefinition} from './content-definition';
import {ErrorUtils} from '../core/error-utils';

@Component({
	selector: 'splbk-create-update-post',
	templateUrl: './create-update-content.component.html',
	styleUrls: ['./create-update-content.component.scss']
})
export class CreateUpdateContentComponent implements OnInit {

	public readonly attachmentParentType: string = 'content';

	private readonly client: ContentClient;

	public definition: ContentDefinition;

	public translate: boolean;
	public postForm: FormGroup;
	public responseError: ResponseError;
	public errorUpload: string;
	public currentCategory: string;

	public utils = UrlUtils;

	// topics
	topics: string[] = [];
	filteredOptionsTopics: Observable<string[]>;
	// tags
	readonly separatorKeysCodes: number[] = [ENTER, COMMA];
	tags: string[] = [];

	languagesNotInserted: string[] = undefined;
	languagesInserted: string[] = [];
	idTranslate: string;

	public content: Content;

	startDate = new Date();
	contentParents: Content[] = [];
	noContentParent = new Content();

	codes: string[] = undefined;
	noCode = '-';
	public languageSelected: string;

	// editor
	@ViewChild('bodyEditor') bodyEditor: any; // CKEditorComponent

	public editor = {
		editor: ClassicEditor, // BaloonEditor, ClassicEditor, BaloonEditorBlocl
		config: {
			language: 'it',
			mediaEmbed: {
				previewsInData: true,
			}
		},
		body: ''
	};
	ckUploadAdapter: CkUploadAdapter;

	public get action(): string {
		return !this.content.isNew() ? 'Update' : 'Create';
	}

	constructor(private spellbook: Spellbook, private route: ActivatedRoute, public dialog: MatDialog, private snackBar: MatSnackBar) {
		this.client = new ContentClient(spellbook);
		this.responseError = undefined;
		this.content = new Content();
	}

	public getCategoryName(): string {
		const cat = this.spellbook.getSupportedCategories().find((s) => {
			return s.name === this.currentCategory;
		});
		return cat ? cat.label : '';
	}

	ngOnInit() {
		this.route.params.subscribe(params => {
			if (params.id) {
				this.content.id = params.id;
				this.client.getContent(this.content.id).subscribe(
					(p: Content) => {
						this.onContentReady(p);
					},
					(error: HttpErrorResponse) => {
						console.error('Error', error);
						this.snackBar.open('Error! ' + error.statusText, 'ok', {});
					});
			}
		});
	}

	private onContentReady(content: Content): void {
		this.content = content;
		this.currentCategory = content.category;
		this.setCurrentType(content.type);
		this.updateForm(content);
		this.setContentParents();

		if (!this.translate) {
			this.setLanguagesNotInserted();
		}
	}

	setTranslateMode(cont: Content, allLanguages: string[], languagesInserted: string[]) {
		this.translate = true;
		// set language
		this.languagesInserted = languagesInserted;

		// set category
		this.idTranslate = cont.idTranslate;
		this.currentCategory = cont.category;
		this.setCurrentType(cont.type);
		// optional insert other Field
		// this.postForm.controls.topic.setValue(cont.topic);

		this.setLanguagesNotInserted(this.languagesInserted, allLanguages);
	}

	public onReadyEditor(data) {
		data.plugins.get('FileRepository').createUploadAdapter = (loader) => {
			this.ckUploadAdapter = new CkUploadAdapter(this.spellbook, 'image', 'content', loader);
			return this.ckUploadAdapter;
		};
	}

	private setContentParents() {
		const filters: Filter[] = [];
		// filters.push(new Filter('Locale', 'it'));
		return this.client.getContentList(0, 999, filters, 'Title', Client.orderAscKey).subscribe((response: ListResponse<Content>) => {
			this.contentParents = response.items;
			this.contentParents.unshift(this.noContentParent);
		});
	}

	private _filterTopics(value: string): string[] {
		const filterValue = value.toLowerCase();
		return this.topics.filter(option => option.toLowerCase().indexOf(filterValue) === 0);
	}

	private setLanguagesNotInserted(supportedLanguage?: string[], allLanguages?: string[]) {
		if (supportedLanguage && allLanguages) {
			// remove all supportedLanguage from allLanguages
			supportedLanguage.forEach((langSupported: string) => {
				const index = allLanguages.findIndex((lang: string) => {
					return lang === langSupported;
				});
				if (index >= 0) {
					// remove it
					allLanguages.splice(index, 1);
					this.languagesNotInserted = allLanguages;
					this.languageSelected = this.languagesNotInserted.length > 0 ? this.languagesNotInserted[0] : null;
					this.postForm.controls.locale.setValue(this.languageSelected); // default
				}
			});
		} else {
			this.spellbook.api.getLanguages().subscribe((allLang: string[]) => {
				if (!this.languagesNotInserted) {
					this.languagesNotInserted = allLang;
					this.languageSelected = this.languagesNotInserted.length > 0 ? this.languagesNotInserted[0] : null;
					this.postForm.controls.locale.setValue(this.languageSelected); // default
				}
			});
		}
	}

	public setCurrentType(type: string) {
		this.content.type = type;
		this.definition = this.spellbook.definitions.getTypeDefinition<Content>(type) as ContentDefinition;
		// if no definition is found, use the default one
		if (!this.definition) {
			this.definition = new ContentDefinition(this.content.type, '', '');
		}

		// now it's safe to build the form
		this.buildForm();
	}

	ValidateDates(group: FormGroup) {
		if (!group || !group.controls || !group.controls.endDate.value) {
			return null;
		}
		if (new Date(group.controls.endDate.value) < (new Date(group.controls.startDate.value))) {
			group.controls.endDate.setErrors({rangeDate: true});
			return {rangeDate: true};
		}
		return null;
	}

	public buildForm(): void {

		this.postForm = new FormGroup({
			title: new FormControl('', [Validators.required, Validators.minLength(2)]),
			locale: new FormControl('', [Validators.required]),
			isPublished: new FormControl(false),
		});

		// now set up defined components
		if (this.definition.field('slug')) {
			this.postForm.addControl('slug', new FormControl(''));
		}
		if (this.definition.field('subtitle')) {
			this.postForm.addControl('subtitle', new FormControl(''));
		}

		if (this.definition.field('topic')) {
			this.postForm.addControl('topic', new FormControl(''));
			// topics
			this.client.getContentProperties('Topic').subscribe((res: string[]) => {
				this.topics = res;
			});
			this.filteredOptionsTopics = this.postForm.controls.topic.valueChanges.pipe(
				startWith(''),
				map(value => this._filterTopics(value))
			);
		}

		if (this.definition.field('order')) {
			this.postForm.addControl('order', new FormControl(''));
			this.postForm.controls.order.setValue(1); // default
		}

		if (this.definition.field('description')) {
			this.postForm.addControl('description', new FormControl(''));
		}

		if (this.definition.field('parent')) {
			this.postForm.addControl('parent', new FormControl(''));
		}

		if (this.definition.field('body')) {
			this.postForm.addControl('body', new FormControl(''));
		}

		if (this.definition.field('code')) {
			this.postForm.addControl('code', new FormControl(''));
			this.client.getSpecialCodeList(this.content.type, this.currentCategory).subscribe((allLang: string[]) => {
				this.codes = allLang;
				this.codes.unshift(this.noCode);
			});
		}

		if (this.definition.field('editor')) {
			this.postForm.addControl('editor', new FormControl(''));
		}

		if (this.definition.field('startDate')) {
			this.postForm.addControl('startDate', new FormControl('', this.definition.isMandatory('startDate') ? [Validators.required] : null));
		}

		if (this.definition.field('endDate')) {
			this.postForm.addControl('endDate', new FormControl('', this.definition.isMandatory('endDate') ? [Validators.required] : null));
		}

		this.definition.addValidationRules(this.postForm);

		this.setContentParents();

		if (!this.translate) {
			this.setLanguagesNotInserted();
		}
	}

	private updateForm(post: Content) {
		this.postForm.patchValue(post);
		this.tags = post.tags;
		if (this.definition.field('startDate') && this.content.startDate) {
			this.startDate = new Date(this.content.startDate);
			this.postForm.controls.startDate.setValue(this.startDate.toISOString());
		}
		if (this.definition.field('endDate')) {
			if (this.content.hasEndDate) {
				this.postForm.controls.endDate.setValue(new Date(this.content.endDate).toISOString());
			} else {
				this.postForm.controls.endDate.setValue(undefined);
			}
		}
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.postForm.controls[controlName].hasError(errorCode);
	}

	public doCreateUpdate(formValue: any): void {
		// this.postForm.controls.endDate.updateValueAndValidity();
		if (this.postForm.valid) {
			if (this.ckUploadAdapter && this.ckUploadAdapter.loadingImage) {
				this.snackBar.open('Wait! Body image uploading...', 'ok', {});
				return;
			}
			// populate the user object
			this.content.title = formValue.title;
			this.content.subtitle = formValue.subtitle;
			this.content.locale = formValue.locale;
			this.content.editor = formValue.editor;
			this.content.revision = formValue.revision;
			this.content.order = formValue.order;
			this.content.description = formValue.description;
			this.content.author = this.spellbook.user.username;
			this.content.tags = this.tags;

			this.content.category = this.currentCategory;
			this.content.topic = formValue.topic;

			this.content.isPublished = formValue.isPublished;
			this.content.body = formValue.body;
			this.content.slug = formValue.slug;
			if (formValue.code !== this.noCode) {
				this.content.code = formValue.code;
			} else {
				this.content.code = '';
			}

			if (formValue.parent) {
				this.content.parent = formValue.parent;
			} else {
				this.content.parent = null;
			}

			this.content.startDate = formValue.startDate ? formValue.startDate : null;
			this.content.endDate = formValue.endDate ? formValue.endDate : null;

			// apply modifications if needed, according to the content definition
			this.definition.beforeSend(this.content);

			if (!this.content.isNew()) {
				this.updatePost();
			} else {
				// create
				if (this.translate) {
					this.content.idTranslate = this.idTranslate;
				}
				this.createPost();
			}
		} else {
			console.error('form not valid', this.content);
			this.responseError = undefined;
			this.responseError = new ResponseError('Error! Form is not valid!', 'generic');
		}
	}

	private createPost(): void {
		this.responseError = undefined;
		this.client.createContent(this.content).subscribe(
			(p: Content) => {
				this.snackBar.open('Content created!', 'ok', {});
				this.spellbook.router.navigate([`/content/${p.id}`]);
			},
			(err: HttpErrorResponse) => {
				this.content.id = undefined;
				this.responseError = ErrorUtils.handlePostError(err, this.postForm);
			}
		);
	}

	private updatePost(): void {
		this.responseError = undefined;
		this.client.updateContent(this.content).subscribe(
			(p: Content) => {
				this.content = p;
				this.updateForm(p);
				this.snackBar.open('Content updated!', 'ok', {});
				// this.snackBar.openFromComponent(SnackbarComponent, {
				// 	data: 'some data'
				// });
			},
			(err: HttpErrorResponse) => {
				this.responseError = ErrorUtils.handlePostError(err, this.postForm);
			}
		);
	}

	delete() {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + this.content.title + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			snackBarRef.dismiss();
			this.client.deleteContent(this.content).subscribe(
				() => {
					this.snackBar.open('Content deleted!', 'ok', {});
					// refresh
					this.spellbook.router.navigate(['/content']);
				},
				(err: HttpErrorResponse) => {
					this.responseError = ErrorUtils.handlePostError(err, this.postForm);
				}
			);

		};
		snackbarData.actionNo = () => {
			snackBarRef.dismiss();
		};
	}

	/**
	 * TAGS
	 */
	add(event: MatChipInputEvent): void {
		const input = event.input;
		const value = event.value;

		// Add our fruit
		if ((value || '').trim()) {
			this.tags.push(value.trim());
		}

		// Reset the input value
		if (input) {
			input.value = '';
		}
	}

	remove(fruit: string): void {
		const index = this.tags.indexOf(fruit);

		if (index >= 0) {
			this.tags.splice(index, 1);
		}
	}

	showAttachmentGroupDialog() {
		const dialogRef = this.dialog.open(CreateAttachmentGroupComponent, {
			width: '500px',
			data: this.content.id,
		});

		dialogRef.afterClosed().subscribe((group: AttachmentGroup) => {
			if (group) {
				this.content.attachmentGroups.unshift(group);
			}
		});
	}

	onCoverClick() {
		let galleryAttachment: SupportedAttachment;
		for (const sa of this.spellbook.supportedAttachments) {
			if (sa.value === AttachmentType.GALLERY) {
				galleryAttachment = sa;
			}
		}
		const dialogRef = this.dialog.open(ListMediaDialogComponent, {
			width: '90vw',
			height: '90vh',
			data: {
				filterableByType: false,
				filterableByGroup: true,
				filterableByParent: false,
				defaultTypeFilter: galleryAttachment,
				defaultGroupFilter: Attachment.GROUP_DEFAULT,
				noCreateAttachment: true
			}
		});
		dialogRef.componentInstance.type = AttachmentType.GALLERY;
		dialogRef.componentInstance.multipleMode = false;
		dialogRef.afterClosed().subscribe((result: Attachment) => {
			if (result) {
				this.content.cover = result.resourceUrl;
			}
		});
	}

	deleteCover() {
		this.content.cover = undefined;
	}

	back() {
		this.spellbook.router.navigate([`../`], { relativeTo: this.route });
	}
}
