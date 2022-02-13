import {Component, Inject, OnInit} from '@angular/core';
import {Spellbook} from '../core/spellbook';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {AttachmentComponent} from './multimedia/attachment.component';
import {Attachment} from './multimedia/attachment';
import {FormControl, FormGroup} from '@angular/forms';
import {SupportedAttachment} from '../core/supported-attachment';
import {FileAttachment} from './multimedia/file';
import {forkJoin} from 'rxjs';
import {HttpErrorResponse} from '@angular/common/http';
import {map} from 'rxjs/operators';

@Component({
	selector: 'splbk-upload-media-file',
	templateUrl: './upload-media-file.component.html',
	styleUrls: ['./upload-media-file.component.css']
})
export class UploadMediaFileComponent implements OnInit {

	type: SupportedAttachment;

	constructor(private spellbook: Spellbook, public dialogRef: MatDialogRef<AttachmentComponent>,
	            @Inject(MAT_DIALOG_DATA) public attachment: Attachment) {
		this.supportedTypes = this.spellbook.supportedAttachments;
		if (attachment) {
			for (const sa of this.spellbook.supportedAttachments) {
				if (sa.value === attachment.type) {
					this.type = sa;
				}
			}
		}

	}

	public multipleMode: boolean;
	public folder: string;
	public namespace: string;
	public errorUpload: string;
	public isLoading: boolean;
	public attachmentForm: FormGroup;
	public supportedTypes: Array<SupportedAttachment>;

	// multiple
	public nUploading: number;
	public nUploadingTot: number;

	static validateAllUrl(group: FormGroup) {
		const url = group.get('url').value;
		const resourceUrl = group.get('resourceUrl').value;
		if (!url && !resourceUrl) {
			return {isRequired: true};
		}
		return null;
	}

	ngOnInit() {
		this.buildForm();
	}

	processFile(imageInput: any) {
		// console.log('files', imageInput.files);
		this.isLoading = true;
		const results = [];
		const fileList: FileList = imageInput.files;
		const self = this;
		this.nUploading = 0;
		this.nUploadingTot = fileList.length;
		Array.from(fileList).forEach((file: File) => {
			const reader = new FileReader();
			this.errorUpload = undefined;

			reader.addEventListener('load', (event: any) => {

				results.push(this.spellbook.api.postFile(
					this.folder ? this.folder : 'image', this.namespace ? this.namespace : 'multimedia', file.name, file
				).pipe(map((f: any) => {
					console.log('f', f);
					this.nUploading++;
					return f;
				})));

				if (results.length === fileList.length) {
					forkJoin(results).subscribe((fileAttachments: FileAttachment[]) => {
							this.isLoading = false;
							fileAttachments.forEach((f: FileAttachment) => {
								f.attachmentType = this.spellbook.getAttachmentTypeForFileAttachment(f);
							});
							if (this.multipleMode) {
								self.dialogRef.close(fileAttachments);
							} else {
								self.dialogRef.close(fileAttachments[0]);
							}
						},
						(error: HttpErrorResponse) => {
							console.error('error upload', error);
							this.isLoading = false;
							if (error.status === 400) {
								this.errorUpload = error.error.Error;
							} else {
								this.errorUpload = error.statusText;
							}
						});
				}
			});

			reader.readAsDataURL(file);
		});

	}

	private buildForm(): void {
		this.attachmentForm = new FormGroup({
			id: new FormControl(''),
			url: new FormControl('', null),
			resourceUrl: new FormControl('', null)
		}, [UploadMediaFileComponent.validateAllUrl]);
	}

	// not used
	public doCreateUpdate(formValue: any): void {
		console.log('doCreateUpdate', formValue);
		if (this.attachmentForm.valid) {
			const result = new FileAttachment();
			result.resourceUrl = formValue.resourceUrl;
			this.dialogRef.close(result);
		}
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.attachmentForm.controls[controlName].hasError(errorCode);
	}

	getFileAcceptForType(type: string) {
		const support = this.supportedTypes.find((s) => {
			return s.value === type;
		});
		return support ? support.accept : '';
	}

}
