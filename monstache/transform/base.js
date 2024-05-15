module.exports = function (doc) {
    return _.pick(
      doc,
      '_id',
      'album',
      'title',
      'artist',
      'genre',
      'release_year',
      'duration'
    );
  };