import {Component, forwardRef, Input, OnInit} from '@angular/core';
import {ControlValueAccessor, NG_VALUE_ACCESSOR} from '@angular/forms';
import {MatSnackBar} from '@angular/material/snack-bar';
import {PlatformLocation} from '@angular/common';

@Component({
	selector: 'splbk-slug',
	template: `
		<mat-form-field fxFlexFill>
			<input matInput type="text" [name]="name" placeholder="{{label}}" [(ngModel)]="slug">
			<button type="button" mat-icon-button matSuffix (click)="makeSlug()"
					[attr.aria-label]="'Create useful slug'">
				<mat-icon>find_replace</mat-icon>
			</button>
			<button type="button" mat-icon-button matSuffix (click)="copyUrl(completeUrl)">
				<mat-icon>file_copy</mat-icon>
			</button>
		</mat-form-field>
	`,
	providers: [
		{
			provide: NG_VALUE_ACCESSOR,
			useExisting: forwardRef(() => SlugComponent),
			multi: true
		}
	],
	styleUrls: ['./slug.component.css']
})
export class SlugComponent implements OnInit, ControlValueAccessor {

	@Input() name: string;
	@Input() label: string;

	@Input()
	src: string;

	@Input()
	intermediateUrl: string;

	completeUrl: string;

	public get baseUrl(): string {
		const url = `${(this.platformLocation as any).location.origin}`;
		return url;
	}

	private _slug: string;

	get slug(): string {
		return this._slug;
	}

	set slug(val: string) {
		this._slug = val;
		if (!this.isAbsoluteUrl()) {
			this.completeUrl = this.baseUrl + '/' + (this.intermediateUrl ? (this.intermediateUrl + '/') : '') + this.slug;
		} else {
			this.completeUrl = this.slug;
		}
		this.propagateChange(this.slug);
	}

	propagateChange = (_: any) => {
	}

	constructor(private snackBar: MatSnackBar, private platformLocation: PlatformLocation) {
	}

	public makeSlug() {
		const src: string = this.src;
		this.slug = src === undefined ? this.slug : src.replace(/[^a-z0-9_]+/gi, '-').replace(/^-|-$/g, '').toLowerCase();
	}

	public copyUrl(value: string) {
		document.addEventListener('copy', (e: ClipboardEvent) => {
			e.clipboardData.setData('text/plain', (value));
			e.preventDefault();
			document.removeEventListener('copy', null);
		});
		document.execCommand('copy');
		this.snackBar.open('\'' + this.completeUrl + '\' copied to the clipboard', 'ok', {});
	}

	isAbsoluteUrl(): boolean {
		const r = new RegExp('^(?:[a-z]+:)?//', 'i');
		return r.test(this.slug);
	}

	ngOnInit(): void {
	}

	registerOnChange(fn: any): void {
		this.propagateChange = fn;
	}

	registerOnTouched(fn: any): void {
	}

	setDisabledState(isDisabled: boolean): void {
	}

	writeValue(obj: any): void {
		if (obj !== undefined) {
			this.slug = obj;
		}
	}
}
