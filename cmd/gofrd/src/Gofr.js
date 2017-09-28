import Ractive from 'ractive';

const Gofr = {};

Gofr.storage = localStorage;

Gofr.helpers = Ractive.defaults.data;

Gofr.helpers.complex = function(r, i) {
  return r + (i < 0 ? "" : " +") + i + "i";
};

Gofr.uuid = function() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

export default Gofr;

