import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Filter} from './client';
import {FormControl, FormGroup, Validators} from '@angular/forms';

@Component({
	selector: 'splbk-search-selector',
	template: `
		<form [formGroup]="formGroup" autocomplete="off" novalidate
				(ngSubmit)="onSubmit(formGroup.value)">
			<div fxLayout="row" fxLayoutAlign="start center" class="my_seacrh">
				<mat-form-field>
					<input matInput placeholder="Filter"
							formControlName="filter"
							(keyup.enter)="onSubmit(formGroup.value)">
					<mat-error *ngIf="hasError('filter', 'required')">Text is required</mat-error>
					<mat-error *ngIf="hasError('filter', 'minLength')">Text is too short</mat-error>
				</mat-form-field>
				<button type="button" color="secondary" mat-raised-button matSuffix (click)="clear()">
					<mat-icon>clear</mat-icon>
				</button>
				<mat-form-field>
					<mat-select placeholder="Field by" formControlName="field">
						<mat-option *ngFor="let field of fields" [value]="field">
							{{field}}
						</mat-option>
					</mat-select>
					<mat-error *ngIf="hasError('field', 'required')">Field is required</mat-error>
				</mat-form-field>
				<button mat-raised-button color="primary" type="submit"
						[disabled]="!formGroup.valid">
					Search
					<mat-icon matListIcon>search</mat-icon>
				</button>
			</div>
		</form>
	`,
	styleUrls: ['./search-selector.component.scss']
})
export class SearchSelectorComponent implements OnInit {

	@Input() fields: string[];

	@Output() filtered: EventEmitter<Filter> = new EventEmitter<Filter>();

	public formGroup: FormGroup;


	ngOnInit(): void {
		this.buildForm();
	}

	private buildForm(): void {
		this.formGroup = new FormGroup({
			filter: new FormControl('', [Validators.required, Validators.minLength(2)]),
			field: new FormControl('', [Validators.required]),
		});
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.formGroup.controls[controlName].hasError(errorCode);
	}

	public onSubmit(formValue: any): void {
		if (this.formGroup.valid) {
			if (formValue.filter && formValue.field) {
				this.filtered.emit(new Filter(formValue.field, formValue.filter));
			}
		}
	}

	clear() {
		this.formGroup.controls.filter.setValue('');
		this.filtered.emit(null);
	}
}
