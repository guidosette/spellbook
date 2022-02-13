import { Injectable } from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ScreenService {

	constructor() {
		this.setSize(new ScreenSize(window.innerWidth));
		window.onresize = () => {
			// set screenWidth on screen size change
			this.setSize(new ScreenSize(window.innerWidth));
		};
	}

	private screenSizeSubject: BehaviorSubject<ScreenSize> = new BehaviorSubject<ScreenSize>(null);

	private setSize(screenSize: ScreenSize): void {
		this.screenSizeSubject.next(screenSize);
	}

	get screenSize(): BehaviorSubject<ScreenSize> {
		return this.screenSizeSubject;
	}
}

export class ScreenSize {
	public width: number;

	constructor(width: number) {
		this.width = width;
	}
}
