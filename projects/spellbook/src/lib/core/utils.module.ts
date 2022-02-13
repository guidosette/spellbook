import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {HttpClientModule} from '@angular/common/http';
import {
	MatButtonModule,
	MatCardModule,
	MatFormFieldModule,
	MatIconModule,
	MatInputModule,
	MatListModule,
	MatMenuModule,
	MatSelectModule,
	MatSidenavModule,
	MatToolbarModule
} from '@angular/material';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {FlexLayoutModule} from '@angular/flex-layout';
import {MenuItemComponent} from './menu-item.component';
import {SnackbarComponent} from './snackbar.component';
import {SlugComponent} from './slug.component';
import {SearchSelectorComponent} from './search-selector.component';
import {RouterModule} from '@angular/router';

@NgModule({
	declarations: [
		MenuItemComponent,
		SnackbarComponent,
		SlugComponent,
		SearchSelectorComponent
	],
	imports: [
		MatCardModule,
		MatInputModule,
		MatFormFieldModule,
		MatButtonModule,
		MatToolbarModule,
		MatSidenavModule,
		MatListModule,
		MatIconModule,
		MatMenuModule,
		MatSelectModule,
		ReactiveFormsModule,
		FormsModule,
		CommonModule,
		HttpClientModule,
		FlexLayoutModule,
		RouterModule
	],
	exports: [
		MenuItemComponent,
		SnackbarComponent,
		SlugComponent,
		SearchSelectorComponent,
	],
})

export class UtilsModule {
}
