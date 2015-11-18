(function() {
  'use-strict';

  xtag.register('saint-laurent-time', {
    shadow: function() {
      /*
          <style type="text/css">
            @import url(https://fonts.googleapis.com/css?family=Roboto+Condensed);

            .saint-laurent-time {
              font-family: 'Roboto Condensed', sans-serif;

              display:inline-block;
              position:relative;
              margin: 0px;
              padding: 0px;

              height: 22px;
              line-height: 22px;
              vertical-align:middle;
            }
          </style>
          <div class="saint-laurent-time" id="container"></div>
      */
    },
    lifecycle: {
      created: function() {
        this.$container = this.shadowRoot.querySelector("#container");
        this._updateTime(this.getAttribute('time'));
      }
    },
    methods: {
      _updateTime: function(time) {
        console.log('update time', time);
        let theDate, display;

        theDate = moment(time);
        display = theDate.format("HH:mm");

        this.$container.innerHTML = display + " (" + theDate.fromNow() + ")";
      },
      accessors: {
        time: {
          attribute: {},
          set: function(value) {
            console.log('set');
            this._updateTime(value);
          }
        }
      }
    }
  });
})();
