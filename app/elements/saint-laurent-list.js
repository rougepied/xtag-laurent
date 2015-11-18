/*
The `saint-laurent-list` web component display a list of next departures
from the Saint-Laurent bus station.

exemple:

    <saint-laurent-list list="[]""></saint-laurent-list>
    <saint-laurent-list list="[{"time":"2015-10-31T22:02:35+01:00","line":"0009"},{"time":"2015-10-31T22:17:00+01:00","line":"0009"}]""></saint-laurent-list>
*/
(function() {

  const template = `
	  <style type="text/css">
		.saint-laurent {
		  display: flex;
		  flex-direction: column;
		}

		.saint-laurent-station {
		  display: inline;
		  position: relative;
		  margin: 0px;
		  padding: 5px;
		}
	  </style>
	  <div class="saint-laurent" id="container"></div>`.trim();

  const sortSchedules = (a, b) => {
    if (moment(a.time).isBefore(b.time)) {
      return -1;
    }
    if (moment(a.time).isAfter(b.time)) {
      return 1;
    }
    return 0;
  }

  const templateStation = (s) => `
      <div class="saint-laurent-station">
        <saint-laurent-badge line="${s.line}"></saint-laurent-badge>
        <saint-laurent-time time="${s.time}"></saint-laurent-time>
      </div>
  `.trim();

  xtag.register('saint-laurent-list', {
    shadow: template,
    lifecycle: {
      created: function() {
        this.$container = this.shadowRoot.querySelector('#container');
      }
    },
    methods: {
      _updateView: function(schedules) {
        let content = "";

        schedules.sort(sortSchedules)
          .forEach(i => {
            content += templateStation(i);
          });

        this.$container.innerHTML = content;
      }
    },
    accessors: {
      list: {
        attribute: {},
        set: function(value) {
          let schedules = JSON.parse(value);
          this._updateView(schedules);
        }
      }
    }
  });
})();
