Ext.define('IW.model.FavoriteWiki', {
    extend: 'Ext.data.Model',
    fields: [
    'Wiki_id', 
    'Username', 
    'Modified', 
    'Status',
    ],
    idgen: 'uuid',
    idProperty: 'Wiki_id',
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
    }
});