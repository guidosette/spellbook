import {AfterViewInit, Component, Inject} from '@angular/core';
import {MAT_SNACK_BAR_DATA} from '@angular/material';

@Component({
	selector: 'splbk-snackbar',
	templateUrl: './snackbar.component.html',
	styleUrls: ['./snackbar.component.css']
})
export class SnackbarComponent {

	constructor(@Inject(MAT_SNACK_BAR_DATA) public snackbarData: SnackbarData) {
	}
}

export class SnackbarData {

	message: string;
	isError: boolean;
	actionOk: () => void;
	actionNo: () => void;
}
