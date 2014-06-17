Ext.define('IW.store.FavoriteWikis', {
    extend: 'Ext.data.Store',
    model: 'IW.model.FavoriteWiki',
    autoLoad: false,
    proxy: {
        type: 'rest',
        url: '/api/favoritewikis',
        reader: {
            type: 'json',
            successProperty: 'success'
        },
        writer: {
            type: 'json'
        }
    },
    listeners: {
            'load' :  function(favorites,records,options) {
                // This is loaded, so fire event to update all the favorite attributes on Wikis
                var wikis = Ext.data.StoreManager.lookup('Wikis');
                wikis.each(function(wiki) {
                    favorites.each(function(favorite) {
                        if (wiki.id==favorite.id) {
                            wiki.set("Favorite", true)
                        } else {
                            wiki.set("Favorite", false)
                        }
                    });
                });
            }
    }
});

