import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';

import {AbstractControl, FormControl, FormGroup, Validators} from '@angular/forms';

import {MatSnackBar} from '@angular/material';
import {HttpErrorResponse} from '@angular/common/http';
import {User} from '../core/user';
import {ResponseError} from '../core/response-error';
import {Spellbook} from '../core/spellbook';
import {ErrorUtils} from '../core/error-utils';

@Component({
	selector: 'splbk-create-update-user',
	templateUrl: './create-update-user.component.html',
	styleUrls: ['./create-update-user.component.css']
})
export class CreateUpdateUserComponent implements OnInit {

	public user: User;
	public userForm: FormGroup;
	public responseError: ResponseError;
	public permissions: PermissionData[];

	public get action(): string {
		return this.user.username ? 'Update' : 'Create';
	}

	constructor(private spellbook: Spellbook, private route: ActivatedRoute, private snackBar: MatSnackBar) {
		this.user = new User();
		this.responseError = undefined;
		this.permissions = [];
		this.permissions.push(new PermissionData('Edit permissions', User.PERMISSION_EDIT_PERMISSIONS));

		const user: PermissionData = new PermissionData('User');
		user.addChildren(new PermissionData('Read', User.PERMISSION_READ_USER));
		user.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_USER));
		this.permissions.push(user);

		const content: PermissionData = new PermissionData('Content');
		content.addChildren(new PermissionData('Read', User.PERMISSION_READ_CONTENT));
		content.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_CONTENT));
		this.permissions.push(content);

		const newsletter: PermissionData = new PermissionData('Newsletter');
		newsletter.addChildren(new PermissionData('Read', User.PERMISSION_READ_MAILMESSAGE));
		newsletter.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_MAILMESSAGE));
		this.permissions.push(newsletter);

		const media: PermissionData = new PermissionData('Media');
		media.addChildren(new PermissionData('Read', User.PERMISSION_READ_MEDIA));
		media.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_MEDIA));
		this.permissions.push(media);

		const place: PermissionData = new PermissionData('Place');
		place.addChildren(new PermissionData('Read', User.PERMISSION_READ_PLACE));
		place.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_PLACE));
		this.permissions.push(place);

		const page: PermissionData = new PermissionData('Page');
		page.addChildren(new PermissionData('Read', User.PERMISSION_READ_PAGE));
		page.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_PAGE));
		this.permissions.push(page);

		const subscription: PermissionData = new PermissionData('Subscription');
		subscription.addChildren(new PermissionData('Read', User.PERMISSION_READ_SUBSCRIPTION));
		subscription.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_SUBSCRIPTION));
		this.permissions.push(subscription);

		const action: PermissionData = new PermissionData('Action');
		action.addChildren(new PermissionData('Read', User.PERMISSION_READ_ACTION));
		action.addChildren(new PermissionData('Write', User.PERMISSION_WRITE_ACTION));
		this.permissions.push(action);
	}

	/** check if an id has been provided.
	 * if so, we are in update state: recover the user
	 */
	ngOnInit() {
		this.buildForm();
		this.route.params.subscribe(params => {
			if (params.username) {
				this.user.username = params.username;
				this.spellbook.api.getUser(this.user.username).subscribe(
					(u: User) => {
						this.user = u;
						this.updateForm(u);
					},
					(error: HttpErrorResponse) => {
						console.error('Error', error);
						this.snackBar.open('Error! ' + error.statusText, 'ok', {});
					});
			}
		});

	}

	private buildForm(): void {
		this.userForm = new FormGroup({
			username: new FormControl('', [Validators.required, Validators.maxLength(32)]),
			password: new FormControl('', [Validators.required, Validators.minLength(User.PASSWORD_MIN_LEN)]),
			email: new FormControl('', [Validators.required, Validators.email]),
			name: new FormControl(''),
			surname: new FormControl(''),
			enabled: new FormControl('')
		});
		this.permissions.forEach((p: PermissionData) => {
			if (!p.isOnlyLabel()) {
				this.userForm.addControl(p.value, new FormControl(''));
			} else {
				this.userForm.addControl(p.label, new FormControl(''));
			}
			p.children.forEach((children: PermissionData) => {
				if (!children.isOnlyLabel()) {
					this.userForm.addControl(children.value, new FormControl(''));
				}
			});
		});
	}

	private updateForm(user: User) {
		if (user) {
			this.userForm.patchValue(user);
			this.userForm.controls.password.setValidators([Validators.minLength(User.PASSWORD_MIN_LEN)]);
			this.userForm.controls.password.updateValueAndValidity();
			const isEnabled: boolean = user.hasPermission(User.PERMISSION_ENABLED);
			this.userForm.controls.enabled.setValue(isEnabled);

			this.permissions.forEach((item: PermissionData) => {
				this.setControlFormPermission(item);
				item.children.forEach((children: PermissionData) => {
					this.setControlFormPermission(children);
				});
			});
		}
	}

	setControlFormPermission(item: PermissionData) {
		if (!item.isOnlyLabel()) {
			this.userForm.controls[item.value].setValue(this.user.hasPermission(item.value));
		} else {
			this.checkControlFormPermissionFirstLevel(item);
		}
	}

	checkControlFormPermissionFirstLevel(item: PermissionData) {
		let checkedCount = 0;
		item.children.forEach((children: PermissionData) => {
			if (this.user.hasPermission(children.value)) {
				checkedCount++;
			}
		});
		this.userForm.controls[item.label].setValue(checkedCount === item.children.length);
	}

	setUserPermissionFromForm(item: PermissionData) {
		if (!item.isOnlyLabel()) {
			if (this.userForm.controls[item.value].value) {
				this.user.permissions.push(item.value);
			}
		}
	}

	changeValueCheckBox(value, item: PermissionData) {
		const checked = !value.path[0].firstChild.checked;
		item.children.forEach((children: PermissionData) => {
			this.userForm.controls[children.value].setValue(checked);
		});
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.userForm.controls[controlName].hasError(errorCode);
	}

	public hasMandatory(controlName: string): boolean {
		const formField = this.userForm.get(controlName);
		if (!formField.validator) {
			return false;
		}
		const validator = formField.validator({} as AbstractControl);
		return (validator && validator.required);
	}

	public doCreateUpdate(formValue: any): void {
		if (this.userForm.valid) {
			// populate the user object
			this.user.email = formValue.email;
			this.user.password = formValue.password;
			this.user.name = formValue.name;
			this.user.surname = formValue.surname;
			this.user.permissions = [];
			if (formValue.enabled) {
				this.user.permissions.push(User.PERMISSION_ENABLED);
			}
			this.permissions.forEach((p: PermissionData) => {
				this.setUserPermissionFromForm(p);
				p.children.forEach((children: PermissionData) => {
					this.setUserPermissionFromForm(children);
				});
			});
			// check if user has username
			if (this.user.username) {
				this.updateUser(this.user);
			} else {
				// create
				this.user.username = this.user.username ? this.user.username : formValue.username;
				this.createUser(this.user);
			}
		}
	}

	private createUser(user: User): void {
		this.responseError = undefined;
		this.spellbook.api.createUser(user).subscribe(
			(u: User) => {
				this.snackBar.open('User created!', 'ok', {});
				this.spellbook.router.navigate([`/users/${u.username}`]);
			},
			(err: HttpErrorResponse) => {
				this.user.username = undefined;
				this.responseError = ErrorUtils.handlePostError(err, this.userForm);
			}
		);
	}

	private updateUser(user: User): void {
		this.responseError = undefined;
		this.spellbook.api.updateUser(user).subscribe(
			(u: User) => {
				this.snackBar.open('User updated!', 'ok', {});
				this.user = u;
				this.updateForm(u);
			},
			(err: HttpErrorResponse) => {
				this.responseError = ErrorUtils.handlePostError(err, this.userForm);
			}
		);
	}

	showPermissions() {
		return this.spellbook.user.hasPermission(User.PERMISSION_EDIT_PERMISSIONS);
	}
}

export class PermissionData {
	label: string;
	value: string;
	children: PermissionData[] = [];

	constructor(label: string, value?: string) {
		this.label = label;
		this.value = value;
	}

	public addChildren(children: PermissionData) {
		this.children.push(children);
	}

	public isOnlyLabel() {
		return this.value === undefined;
	}
}
