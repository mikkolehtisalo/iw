Ext.define('IW.view.page.PageTree', {
    extend: 'Ext.tree.Panel',
    alias: 'widget.pagetree',
    requires: [
        'Ext.tree.*',
        'Ext.data.*'
    ],
    
    rootVisible: false,
    layout: 'fit',
    border: 0,

    enableDD: true,
    invalidateScrollerOnRefresh: false,

    viewConfig: {
        loadMask: false,
        plugins: {
            ptype: 'treeviewdragdrop',
            containerScroll: true
        },
        listeners: {       
            drop: function (node, data, overModel, dropPosition) {
                // Synchronize all the changes to the server side after drag & drop event
                var moved = data.records;
                for (index in moved) {
                    var record = moved[index];
                    // In case this was not an order change but move, fix paths. The server side will have to work on that.
                    var parentPath = record.parentNode.data.Path;
                    if (parentPath.length > 0) {
                        record.data.Path = parentPath + '/' + record.data.Page_id;
                    } else {
                        // Probably moved to root...
                        record.data.Path = record.data.Page_id;
                    }
                    for (index in record.stores) {
                        var st = record.stores[index];
                        st.sync({
                            success: function() {
                                // Nothing yet I guess
                            },
                            failure: function() {
                                console.log('Unable to save changes');
                            }
                        });
                    }
                }
                        
            },         
        }   
    },

    store: 'Pages',
    displayField: 'Title',
    selModel: { allowDeselect: true },

    initComponent: function() {

        // Creates a new instance of store
        // This is probably slightly unusual, but we will be using the stores a lot later on
        var theStore = Ext.create('IW.store.Pages');
        Ext.apply(this, {
            store: theStore
        }); 

        // Change the root depending on the wiki
        this.getStore().setRootNode({
            expanded: true,
            //id: this.record.data.wiki_id
            Wiki_id: this.record.data.Wiki_id
        });

        this.callParent(arguments);

    }
});
