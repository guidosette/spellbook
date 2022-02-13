import {Component, OnInit} from '@angular/core';
import {Page} from './page';
import {AbstractControl, FormControl, FormGroup, Validators} from '@angular/forms';


import {MatSnackBar} from '@angular/material';
import {ActivatedRoute} from '@angular/router';
import {PlatformLocation} from '@angular/common';
import {HttpErrorResponse} from '@angular/common/http';
import {NavigationClient} from './navigation-client';
import {ResponseError} from '../core/response-error';
import {Spellbook} from '../core/spellbook';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {ErrorUtils} from '../core/error-utils';

@Component({
	selector: 'splbk-create-update-page',
	templateUrl: './create-update-page.component.html',
	styleUrls: ['./create-update-page.component.scss']
})
export class CreateUpdatePageComponent implements OnInit {
	private readonly client: NavigationClient;

	public page: Page;
	public postForm: FormGroup;
	public responseError: ResponseError;

	languages: string[] = undefined;
	codes: string[] = undefined;

	public get baseUrl(): string {
		const url = `${(this.platformLocation as any).location.origin}`;
		return this.page.locale ? `${url}/${this.page.locale}` : url;
	}

	public get action(): string {
		return this.page.id ? 'Update' : 'Create';
	}

	public get canSetUrl(): boolean {
		return !this.page.id;
	}

	constructor(private spellbook: Spellbook, private route: ActivatedRoute, private snackBar: MatSnackBar, private platformLocation: PlatformLocation) {
		this.client = new NavigationClient(spellbook);
		this.responseError = undefined;
		this.page = new Page();
	}

	ngOnInit() {
		this.buildForm();

		this.route.params.subscribe(params => {
			if (params.id) {
				this.page.id = params.id;
				this.client.getPage(this.page.id).subscribe(
					(p: Page) => {
						this.page = p;
						this.updateForm(p);
					},
					(error) => {
						console.error('Error', error);
						this.snackBar.open('Error! ' + error.statusText, 'ok', {});
					});
			}
		});
		this.setLanguages();
		this.setCodes();
	}

	private setLanguages() {
		this.spellbook.api.getLanguages().subscribe((allLang: string[]) => {
			this.languages = allLang;
		});
	}

	private setCodes() {
		this.client.getStaticPageCodeList().subscribe((allLang: string[]) => {
			this.codes = allLang;
		});
	}

	public doCreateUpdate(formValue: any): void {
		if (this.postForm.valid) {
			// populate the user object
			this.page.label = formValue.label;
			this.page.title = formValue.title;
			this.page.metadesc = formValue.metadesc;
			this.page.url = formValue.url;
			this.page.locale = formValue.locale;
			this.page.order = formValue.order;
			this.page.code = formValue.code;

			if (this.page.id) {
				this.updatePage();
			} else {
				this.createPage();
			}
		}
	}

	private createPage(): void {
		this.responseError = undefined;
		this.client.createPage(this.page).subscribe(
			(p: Page) => {
				this.snackBar.open('Page successfully created', 'ok', {});
				this.spellbook.router.navigate([`/page/${p.id}`]);
			},
			(err: HttpErrorResponse) => {
				this.page.id = undefined;
				this.responseError = ErrorUtils.handlePostError(err, this.postForm);
			}
		);
	}

	private updatePage(): void {
		this.responseError = undefined;
		this.client.updatePage(this.page).subscribe(
			(p: Page) => {
				this.page = p;
				this.updateForm(p);
				this.snackBar.open('Page updated!', 'ok', {});
			},
			(err: HttpErrorResponse) => {
				this.responseError = ErrorUtils.handlePostError(err, this.postForm);
			}
		);
	}

	delete() {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + this.page.title + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			console.log('actionOk');
			snackBarRef.dismiss();
			this.client.deletePage(this.page).subscribe(
				() => {
					this.snackBar.open('Page deleted!', 'ok', {});
					// refresh
					this.spellbook.router.navigate(['/page']);
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

	public hasError(controlName: string, errorCode: string): boolean {
		return this.postForm.controls[controlName].hasError(errorCode);
	}

	public hasMandatory(controlName: string): boolean {
		const formField = this.postForm.get(controlName);
		if (!formField.validator) {
			return false;
		}
		const validator = formField.validator({} as AbstractControl);
		return (validator && validator.required);
	}

	private buildForm(): void {
		this.postForm = new FormGroup({
			label: new FormControl('', [Validators.required, Validators.minLength(2)]),
			title: new FormControl(''),
			metadesc: new FormControl(''),
			order: new FormControl(''),
			url: new FormControl(''),
			locale: new FormControl('', [Validators.required]),
			code: new FormControl('', [Validators.required]),
		});
		this.postForm.controls.order.setValue(1); // default
	}

	private updateForm(page: Page) {
		this.postForm.patchValue(page);
	}

}
