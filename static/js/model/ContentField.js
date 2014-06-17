Ext.define('IW.model.ContentField', {
    extend: 'Ext.data.Model',
    fields: [
        'Contentfield_id', 
        'Wiki_id', 
        'Content', 
        'Contentwithmacros'
    ],
    proxy: {
        type: 'rest',
        url: '/api/contentfields/',
        reader: {
            type: 'json',
            successProperty: 'success'
        }
    }
});

