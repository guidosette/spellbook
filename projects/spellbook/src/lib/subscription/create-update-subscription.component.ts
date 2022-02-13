import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';


import {ActivatedRoute} from '@angular/router';
import {MatSnackBar} from '@angular/material';

import {Subscription} from './subscription';
import {HttpErrorResponse} from '@angular/common/http';
import {SubscriptionClient} from './subscription-client';
import {ResponseError} from '../core/response-error';
import {Spellbook} from '../core/spellbook';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {ErrorUtils} from '../core/error-utils';

@Component({
	selector: 'splbk-create-update-newsletter',
	templateUrl: './create-update-subscription.component.html',
	styleUrls: ['./create-update-subscription.component.scss']
})
export class CreateUpdateSubscriptionComponent implements OnInit {
	private readonly client: SubscriptionClient;

	public subscribe: Subscription;
	public postForm: FormGroup;
	public responseError: ResponseError;

	public get action(): string {
		return !this.subscribe.isNew() ? 'Update' : 'Create';
	}

	constructor(private spellbook: Spellbook, private route: ActivatedRoute, private snackBar: MatSnackBar) {
		this.responseError = undefined;
		this.subscribe = new Subscription();
		this.client = new SubscriptionClient(spellbook);
	}

	ngOnInit() {
		this.buildForm();

		this.route.params.subscribe(params => {
			if (params.key) {
				this.subscribe.key = params.key;
				this.client.getSubscription(this.subscribe.key).subscribe(
					(m: Subscription) => {
						this.subscribe = m;
						this.updateForm(m);
					},
					(error) => {
						console.error('Error', error);
						this.snackBar.open('Error! ' + error.statusText, 'ok', {});
					});
			}
		});
	}

	delete() {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + this.subscribe.email + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			console.log('actionOk');
			snackBarRef.dismiss();
			this.client.deleteSubscription(this.subscribe).subscribe(
				() => {
					this.snackBar.open('Subscribe deleted!', 'ok', {});
					// refresh
					this.spellbook.router.navigate(['/subscription']);
				},
				(err: HttpErrorResponse) => {
					this.responseError = ErrorUtils.handlePostError(err, this.postForm);
				}
			);

		};
		snackbarData.actionNo = () => {
			console.log('actionNo');
			snackBarRef.dismiss();
		};
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.postForm.controls[controlName].hasError(errorCode);
	}

	private buildForm(): void {
		this.postForm = new FormGroup({
			email: new FormControl('', [Validators.required, Validators.minLength(2)]),
			country: new FormControl('', [Validators.required, Validators.minLength(2)]),
			firstName: new FormControl('', [Validators.required, Validators.minLength(2)]),
			lastName: new FormControl('', [Validators.required, Validators.minLength(2)]),
			organization: new FormControl('', [Validators.required, Validators.minLength(2)]),
			position: new FormControl('', [Validators.required, Validators.minLength(2)]),
			notes: new FormControl(''),
		});
	}

	private updateForm(m: Subscription) {
		this.postForm.patchValue(m);
	}

}
