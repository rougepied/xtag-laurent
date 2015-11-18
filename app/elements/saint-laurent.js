(function() {
  'use-strict';

  const newStationParam = (stop, route, direction) => [
      ["stop", stop],
      ["route", route],
      ["direction", direction]
    ].map(i => i.join("="));



  xtag.register('saint-laurent', {
    content: '<saint-laurent-list id="saint-laurent-list" list="[]"/>',
    lifecycle: {
      created: function() {
        this.stations = [
          newStationParam("1372", "0009", "0"),
          newStationParam("1485", "0071", "0"),
          newStationParam("1103", "0003", "0"),
          newStationParam("1485", "0051", "0")
        ];

        let results = [];
        this.stations.forEach(i => {
          let param = i.join("&");
          console.log(param);

          fetch("/api/3.0?" + param)
            .then(i => i.json())
            .then(i => {
              if (i !== null) {
                results = results.concat(i.schedules);
                this._updateList(JSON.stringify(results));
              }
            });
        });
      }
    },
    methods: {
      _updateList: function(list) {
        console.log('update list');
        this.querySelector('saint-laurent-list').setAttribute('list', list);
      }
    }
  });
})();
