import {Component, Inject, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';

import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {CreateUpdateMediaComponent} from './create-update-media.component';
import {FileAttachment} from './multimedia/file';
import {SupportedAttachment} from '../core/supported-attachment';
import {Attachment} from './multimedia/attachment';
import {Spellbook} from '../core/spellbook';
import {UrlUtils} from '../core/url-utils';

@Component({
	selector: 'splbk-collect-media-url',
	templateUrl: './collect-media-url.component.html',
	styleUrls: ['./collect-media-url.component.css']
})
export class CollectMediaUrlComponent implements OnInit {
	public isLoading: boolean;
	public attachmentForm: FormGroup;
	public errorMalformedUrl: boolean;
	public errorWrongUrl: boolean;
	public errorWrongMime: boolean;
	type: SupportedAttachment;

	public utils = UrlUtils;

	constructor(private spellbook: Spellbook, public dialogRef: MatDialogRef<CreateUpdateMediaComponent>, @Inject(MAT_DIALOG_DATA) atype: Attachment) {
		if (atype) {
			for (const sa of this.spellbook.supportedAttachments) {
				if (sa.value === atype.type) {
					this.type = sa;
				}
			}
		}
	}

	ngOnInit() {
		this.buildForm();
	}

	private buildForm(): void {
		this.attachmentForm = new FormGroup({
			url: new FormControl('', [Validators.required, Validators.minLength(2)]),
		});

		this.attachmentForm.get('url').valueChanges.subscribe((value: string) => {
			if (!value) {
				return;
			}
			this.errorMalformedUrl = false;
			this.errorWrongUrl = false;
			this.errorWrongMime = false;
		});
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.attachmentForm.controls[controlName].hasError(errorCode);
	}

	public submit(formValue: any) {
		if (this.attachmentForm.valid) {
			this.errorMalformedUrl = false;
			this.errorWrongUrl = false;
			this.errorWrongMime = false;

			const url = formValue.url;
			if (url.startsWith('http://') || url.startsWith('https://')) {

				const result = new FileAttachment();
				result.resourceUrl = url;
				result.attachmentType = this.spellbook.getAttachmentTypeForFileAttachment(result);
				this.dialogRef.close(result);
			} else {
				this.errorMalformedUrl = true;
			}
		}
	}

	public close() {
		this.dialogRef.close();
	}

}
