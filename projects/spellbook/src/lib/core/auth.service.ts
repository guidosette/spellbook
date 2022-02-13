import {Injectable} from '@angular/core';
import {ActivatedRouteSnapshot, CanActivate, CanActivateChild, RouterStateSnapshot} from '@angular/router';
import {Observable} from 'rxjs/internal/Observable';
import {Spellbook} from './spellbook';

@Injectable()
export class AuthService implements CanActivate, CanActivateChild {

	constructor(private spellbook: Spellbook) {

	}

	canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
		if (this.spellbook.isLoggedIn()) {
			return true;
		} else {
			if (this.spellbook.router) {
				this.spellbook.router.navigate([`/login`]);
			}
			return false;
		}
	}

	canActivateChild(childRoute: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
		return this.canActivate(childRoute, state);
	}

}
