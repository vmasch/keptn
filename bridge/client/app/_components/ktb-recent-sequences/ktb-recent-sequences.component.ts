import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input} from '@angular/core';
import {Project} from '../../_models/project';
import {DateUtil} from '../../_utils/date.utils';

@Component({
  selector: 'ktb-recent-sequences',
  templateUrl: './ktb-recent-sequences.component.html',
  styleUrls: ['./ktb-recent-sequences.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class KtbRecentSequencesComponent {
  private _project: Project;

  @Input()
  get project() {
    return this._project;
  }
  set project(project) {
    if (this._project !== project || this._project.sequences !== project.sequences) {
      this._project = project;
      this._changeDetectorRef.markForCheck();
    }
  }
  constructor(public dateUtil: DateUtil, private _changeDetectorRef: ChangeDetectorRef) { }

}
