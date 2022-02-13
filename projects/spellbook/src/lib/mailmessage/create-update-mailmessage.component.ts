import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';


import {ActivatedRoute} from '@angular/router';
import {MatSnackBar} from '@angular/material';

import {MailMessage} from './mailmessage';
import {HttpErrorResponse} from '@angular/common/http';
import {MessageClient} from './message-client';
import {ResponseError} from '../core/response-error';
import {Spellbook} from '../core/spellbook';
import {SnackbarComponent, SnackbarData} from '../core/snackbar.component';
import {ErrorUtils} from '../core/error-utils';

@Component({
	selector: 'splbk-create-update-mailmessage',
	templateUrl: './create-update-mailmessage.component.html',
	styleUrls: ['./create-update-mailmessage.component.scss']
})
export class CreateUpdateMailMessageComponent implements OnInit {
	private readonly client: MessageClient;

	public mailMessage: MailMessage;
	public postForm: FormGroup;
	public responseError: ResponseError;

	public get action(): string {
		return this.mailMessage.id ? 'Update' : 'Create';
	}

	constructor(private spellbook: Spellbook, private route: ActivatedRoute, private snackBar: MatSnackBar) {
		this.responseError = undefined;
		this.mailMessage = new MailMessage();
		this.client = new MessageClient(spellbook);
	}

	ngOnInit() {
		this.buildForm();

		this.route.params.subscribe(params => {
			if (params.id) {
				this.mailMessage.id = params.id;
				this.client.getMailMessage(this.mailMessage.id).subscribe(
					(m: MailMessage) => {
						this.mailMessage = m;
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
		snackbarData.message = 'Are you sure to delete ' + this.mailMessage.object + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			console.log('actionOk');
			snackBarRef.dismiss();
			this.client.deleteMailMessage(this.mailMessage).subscribe(
				() => {
					this.snackBar.open('Mail Message deleted!', 'ok', {});
					// refresh
					this.spellbook.router.navigate(['/mailmessage']);
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
			sender: new FormControl('', [Validators.required, Validators.minLength(2)]),
			recipient: new FormControl('', [Validators.required, Validators.minLength(2)]),
			object: new FormControl(''),
			body: new FormControl('', [Validators.required, Validators.minLength(2)]),
		});
	}

	private updateForm(m: MailMessage) {
		this.postForm.patchValue(m);
	}

}
