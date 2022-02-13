import {Component, ElementRef, Inject, Input, OnInit, ViewChild} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef, MatSnackBar} from '@angular/material';
import {Attachment} from './attachment';
import {FormControl, FormGroup, Validators} from '@angular/forms';


import {HttpErrorResponse} from '@angular/common/http';
import {SupportedAttachment} from '../../core/supported-attachment';

import {FileAttachment} from './file';
import {MediaClient} from '../media-client';
import {ResponseError} from '../../core/response-error';
import {Spellbook} from '../../core/spellbook';
import {SnackbarComponent, SnackbarData} from '../../core/snackbar.component';
import {ListResponse} from '../../core/client';
import {ErrorUtils} from '../../core/error-utils';

@Component({
	selector: 'splbk-multimedia',
	templateUrl: './attachment.component.html',
	styleUrls: ['./attachment.component.css']
})
export class AttachmentComponent implements OnInit {
	private readonly client: MediaClient;

	@Input() id: string;
	@ViewChild('resource')
	inputImage: ElementRef;

	public attachmentForm: FormGroup;
	public responseError: ResponseError;
	public errorUpload: string;
	public isLoading: boolean;

	urlFiles: FileAttachment[] = [];
	selectedValue: string;

	public supportedTypes: Array<SupportedAttachment>;

	public get action(): string {
		return this.id ? 'Update' : 'Create';
	}

	constructor(private spellbook: Spellbook,
	            public dialogRef: MatDialogRef<AttachmentComponent>,
	            @Inject(MAT_DIALOG_DATA) public attachment: Attachment,
	            private snackBar: MatSnackBar) {

		this.client = new MediaClient(spellbook);
		this.responseError = undefined;
		this.supportedTypes = this.spellbook.supportedAttachments;

		this.setUrlFiles();
	}

	private setUrlFiles() {
		this.spellbook.api.getUrlFiles().subscribe((res: ListResponse<FileAttachment>) => {
			this.urlFiles = res.items;
		});
	}

	getImageForType(type: string) {
		const support = this.supportedTypes.find((s) => {
			return s.value === type;
		});
		return support ? support.image : '';
	}

	getNameForType(type: string) {
		const support = this.supportedTypes.find((s) => {
			return s.value === type;
		});
		return support ? support.name : '';
	}

	getFileAcceptForType(type: string) {
		const support = this.supportedTypes.find((s) => {
			return s.value === type;
		});
		return support ? support.accept : '';
	}

	ngOnInit() {
		this.buildForm();

		if (this.attachment.id) {
			this.id = this.attachment.id;
			this.updateForm(this.attachment);
		} else if (this.attachment.group) {
			this.updateForm(this.attachment);
		}
	}

	processFile(imageInput: any) {
		this.isLoading = true;
		const file: File = imageInput.files[0];
		const reader = new FileReader();
		this.errorUpload = undefined;

		reader.addEventListener('load', (event: any) => {
			this.spellbook.api.postFile(
				'image', 'multimedia', file.name, file
			).subscribe((res: FileAttachment) => {
					this.isLoading = false;
					this.attachmentForm.controls.resourceUrl.setValue(res.resourceUrl);
					this.attachmentForm.controls.name.setValue(file.name);
				},
				(error: HttpErrorResponse) => {
					this.isLoading = false;
					this.errorUpload = error.statusText;
				});
		});

		reader.readAsDataURL(file);
	}

	private buildForm(): void {
		this.attachmentForm = new FormGroup({
			id: new FormControl(''),
			parentKey: new FormControl('', [Validators.required]),
			name: new FormControl('', [Validators.required]),
			description: new FormControl(''),
			type: new FormControl('', [Validators.required]),
			group: new FormControl('', [Validators.required]),
			url: new FormControl('', null),
			manualUrl: new FormControl('', null),
			selectUrl: new FormControl('', null),
			// selected: new FormControl('', null),
			resourceUrl: new FormControl('', null)
		}, [this.validateAllUrl]);

		this.attachmentForm.get('selectUrl').valueChanges.subscribe((value: FileAttachment) => {
			if (!value) {
				return;
			}
			this.attachmentForm.controls.resourceUrl.patchValue(value.resourceUrl);
			this.attachmentForm.controls.name.setValue(value.name);
		});

		this.attachmentForm.get('manualUrl').valueChanges.subscribe((value: string) => {
			if (!value) {
				return;
			}
			this.attachmentForm.controls.resourceUrl.patchValue(value);
			this.attachmentForm.controls.name.setValue(value);
		});
	}

	clickNewFile() {
		this.selectedValue = 'file';
		this.attachmentForm.controls.selectUrl.patchValue(undefined); // unselect
		this.attachmentForm.controls.resourceUrl.patchValue(undefined); // unselect
	}

	clickListFile() {
		this.selectedValue = 'select';
		// unselect url input
		this.attachmentForm.controls.url.patchValue(undefined);
		if (this.inputImage) {
			this.inputImage.nativeElement.value = '';
		}
		this.attachmentForm.controls.resourceUrl.patchValue(undefined); // unselect
	}

	clickInsertManualUrl() {
		this.selectedValue = 'manualUrl';
		// unselect url input
		this.attachmentForm.controls.url.patchValue(undefined);
		if (this.inputImage) {
			this.inputImage.nativeElement.value = '';
		}
		this.attachmentForm.controls.resourceUrl.patchValue(undefined); // unselect
	}

	validateAllUrl(group: FormGroup) {
		const url = group.get('url').value;
		const selectUrl = group.get('selectUrl').value;
		const resourceUrl = group.get('resourceUrl').value;
		if (!url && !selectUrl && !resourceUrl) {
			return {isRequired: true};
		}
		return null;
	}

	private updateForm(attachment: Attachment) {
		this.attachmentForm.patchValue(attachment);
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.attachmentForm.controls[controlName].hasError(errorCode);
	}

	public doCreateUpdate(formValue: any): void {
		if (this.attachmentForm.valid) {
			// populate the user object
			const attachment = new Attachment();
			attachment.id = this.id ? this.id : formValue.id;
			attachment.parentKey = formValue.parentKey;
			attachment.name = formValue.name;
			attachment.type = formValue.type;
			attachment.description = formValue.description;
			attachment.group = formValue.group;
			attachment.resourceUrl = formValue.resourceUrl;

			// check if slug has group
			if (this.id) {
				this.updateAttachment(attachment);
			} else {
				// create
				this.createAttachment(attachment);
			}
		}
	}

	private createAttachment(attachment: Attachment): void {
		this.responseError = undefined;
		this.client.createAttachment(attachment).subscribe(
			(m: Attachment) => {
				this.updateForm(m);
				this.dialogRef.close(m);
			},
			(err: HttpErrorResponse) => {
				this.responseError = ErrorUtils.handlePostError(err, this.attachmentForm);
			}
		);
	}

	private updateAttachment(attachment: Attachment): void {
		this.responseError = undefined;
		this.client.updateAttachment(attachment).subscribe(
			(m: Attachment) => {
				this.updateForm(m);
				this.dialogRef.close(m);
			},
			(err: HttpErrorResponse) => {
				this.responseError = ErrorUtils.handlePostError(err, this.attachmentForm);
			}
		);
	}

	delete() {
		const snackbarData: SnackbarData = new SnackbarData();
		snackbarData.message = 'Are you sure to delete ' + this.attachment.name + '?';
		const snackBarRef = this.snackBar.openFromComponent(SnackbarComponent, {
			duration: 30000,
			data: snackbarData
		});
		snackbarData.actionOk = () => {
			console.log('actionOk');
			snackBarRef.dismiss();
			this.client.deleteAttachment(this.attachment).subscribe(
				() => {
					this.snackBar.open('Attachment deleted!', 'ok', {});
					this.attachment.id = undefined;
					this.dialogRef.close(this.attachment);
				},
				(err: HttpErrorResponse) => {
					this.responseError = ErrorUtils.handlePostError(err, this.attachmentForm);
				}
			);

		};
		snackbarData.actionNo = () => {
			console.log('actionNo');
			snackBarRef.dismiss();
		};
	}
}
