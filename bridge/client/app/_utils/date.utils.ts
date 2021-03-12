import * as moment from "moment";
import {Trace} from "../_models/trace";
import {Injectable} from "@angular/core";

@Injectable({
  providedIn: 'root'
})
export class DateUtil {

  public DEFAULT_DATE_FORMAT = 'YYYY-MM-DD';
  public DEFAULT_TIME_FORMAT = 'HH:mm';

  public getDurationFormatted(start, end?) {
    const duration = this.getDuration(start, end);

    let result = duration.seconds+' seconds';
    if(duration.minutes > 0)
      result = duration.minutes+' minutes '+result;
    if(duration.hours > 0)
      result = duration.hours+' hours '+result;
    if(duration.days > 0)
      result = duration.days+' days '+result;

    return result;
  }

  private getDuration(start, end?) {
    const diff = moment(end).diff(moment(start));
    const duration = moment.duration(diff);

    const days = Math.floor(duration.asDays());
    const hours = Math.floor(duration.asHours()%24);
    const minutes = Math.floor(duration.asMinutes()%60);
    const seconds = Math.floor(duration.asSeconds()%60);
    return {days, hours, minutes, seconds};
  }

  public getDurationFormattedShort(start, end?) {
    const duration = this.getDuration(start, end);

    let result = '';
    if (duration.days > 0) {
      result = duration.days + ' day' + (duration.days === 1 ? '' : 's');
    } else if (duration.hours > 0) {
      result = duration.hours + ' hour' + (duration.hours === 1 ? '' : 's');
    } else if (duration.minutes > 0) {
      result = duration.minutes + ' min.';
    } else {
      result = duration.seconds + ' sec.';
    }

    return result + ' ago';
  }

  public getCalendarFormats(showSeconds?: boolean) {
    if(showSeconds) {
      return {
        lastDay : '[yesterday at] HH:mm:ss',
        sameDay : '[today at] HH:mm:ss',
        nextDay : '[tomorrow at] HH:mm:ss',
        lastWeek : '[last] dddd [at] HH:mm:ss',
        nextWeek : 'dddd [at] HH:mm:ss',
        sameElse : 'YYYY-MM-DD HH:mm:ss'
      };
    }
    return {
      lastDay : '[yesterday at] HH:mm',
      sameDay : '[today at] HH:mm',
      nextDay : '[tomorrow at] HH:mm',
      lastWeek : '[last] dddd [at] HH:mm',
      nextWeek : 'dddd [at] HH:mm',
      sameElse : 'YYYY-MM-DD HH:mm'
    };
  }

  public getDateTimeFormat() {
    return [this.DEFAULT_DATE_FORMAT, this.DEFAULT_TIME_FORMAT].join(" ");
  }

  public getTimeFormat() {
    return this.DEFAULT_TIME_FORMAT;
  }

  static compareTraceTimesAsc(a: Trace, b: Trace) {
    return new Date(b.time).getTime() - new Date(a.time).getTime();
  }

  static compareTraceTimesDesc(a: Trace, b: Trace) {
    return new Date(a.time).getTime() - new Date(b.time).getTime();
  }
}
