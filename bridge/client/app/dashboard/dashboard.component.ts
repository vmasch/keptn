import {ChangeDetectorRef, Component, OnDestroy, OnInit} from '@angular/core';
import {Observable, Subject, timer} from 'rxjs';

import {Project} from '../_models/project';

import {DataService} from '../_services/data.service';
import {takeUntil, filter} from 'rxjs/operators';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy{

  public projects: Project[];
  private _recentSequenceTimerInterval = 30 * 1000; // milliseconds
  private readonly unsubscribe$ = new Subject<void>();

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
    this.dataService.projects
      .pipe(
        takeUntil(this.unsubscribe$),
        filter(projects => !!projects)
      ).subscribe(() => {
        this.loadProjects();
      });
  }

  ngOnInit(): void {
    timer(this._recentSequenceTimerInterval, this._recentSequenceTimerInterval)
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(() => {
        this.loadProjects();
      });
  }

  loadProjects() {
    this.dataService.loadRecentSequences()
      .subscribe(projects => {
        this.projects = projects;
      });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
