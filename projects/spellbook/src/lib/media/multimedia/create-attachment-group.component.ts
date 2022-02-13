import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {AttachmentGroup} from './attachment-group';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {Spellbook} from '../../core/spellbook';
import {SupportedAttachment} from '../../core/supported-attachment';


/**
 * Created by Luigi Tanzini (luigi.tanzini@distudioapp.com) on 2019-02-28.
 */
@Component({
	selector: 'splbk-create-gallery',
	templateUrl: './create-attachment-group.component.html',
	styleUrls: ['./create-attachment-group.component.css']
})
export class CreateAttachmentGroupComponent implements OnInit {

	public supportedTypes: Array<SupportedAttachment>;

	public groupForm: FormGroup;
	private readonly parentKey: string;

	constructor(private spellbook: Spellbook,
	            public dialogRef: MatDialogRef<CreateAttachmentGroupComponent>,
	            @Inject(MAT_DIALOG_DATA) data: string) {
		this.parentKey = data; // encode key of content
		this.supportedTypes = this.spellbook.supportedAttachments;
	}

	ngOnInit(): void {
		this.buildForm();
	}

	private buildForm(): void {
		this.groupForm = new FormGroup({
			type: new FormControl('', [Validators.required]),
			name: new FormControl('', [Validators.minLength(3), Validators.required]),
		});
	}

	public createAttachmentGroup(formValue: any) {
		if (this.groupForm.valid) {
			// populate the user object
			const type = formValue.type;
			const name = formValue.name;
			const group = new AttachmentGroup(name, type, this.parentKey);
			this.dialogRef.close(group);
		}
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.groupForm.controls[controlName].hasError(errorCode);
	}
}
