import {Component, OnInit} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';
import {Spellbook} from './spellbook';
import {MatSnackBar} from '@angular/material';

@Component({
	selector: 'splbk-login',
	templateUrl: './login.component.html',
	styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

	public loginForm: FormGroup;
	public responseError: string;

	constructor(private spellbook: Spellbook, private snackBar: MatSnackBar) {
		if (this.spellbook.isLoggedIn()) {
			// user alredy loggin
			if (this.spellbook.router) {
				this.spellbook.router.navigate([`/`]);
			}
		}
	}

	ngOnInit() {
		this.prepareForm();
	}

	private prepareForm(): void {
		this.loginForm = new FormGroup({
			username: new FormControl('', [Validators.required, Validators.maxLength(32)]),
			password: new FormControl('', [Validators.required, Validators.minLength(8)]),
		});
	}

	public hasError(controlName: string, errorCode: string): boolean {
		return this.loginForm.controls[controlName].hasError(errorCode);
	}

	public doLogin(formValue: any): void {
		if (this.loginForm.valid) {
			this.login(formValue.username, formValue.password);
		}
	}

	private login(username: string, password: string): void {
		this.responseError = '';
		this.spellbook.api.createToken(username, password).subscribe(
			(token: string) => {
				this.spellbook.token = token;
				// navigate to the dashboard
				this.spellbook.init().finally(()=>{
					this.spellbook.router.navigate(['/']);
				});
			},
			(err) => {
				this.responseError = err.statusText;
				this.snackBar.open('Error ' + err.statusText, 'ok', {});
			}
		);
	}

	loginGoogle() {
		// console.log('loginGoogle superUserUrl', this.spellbook.superUserUrl);
		// this.spellbook.router.navigate([this.spellbook.superUserUrl]);
		window.location.href = this.spellbook.superUserUrl;
	}
}
