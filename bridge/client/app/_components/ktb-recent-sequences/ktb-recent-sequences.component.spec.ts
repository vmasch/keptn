import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbRecentSequencesComponent } from './ktb-recent-sequences.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbRecentSequencesComponent', () => {
  let component: KtbRecentSequencesComponent;
  let fixture: ComponentFixture<KtbRecentSequencesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbRecentSequencesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
