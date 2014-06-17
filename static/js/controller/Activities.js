Ext.define('IW.controller.Activities', {
    extend: 'Ext.app.Controller',

    views: [
    'activity.Window',
    'activity.List'
    ],
    stores: [
    'Activities', 
    'Wikis',
    'Pages'
    ],
    models: [
    'Activity',
    ],
    init: function() {
        this.control({
            'activitylist': {
                itemdblclick: this.openTarget,
                itemdeleteclick: this.deleteWiki
            }
        });
    },
    openTarget: function(grid, record, item, index, e, eOpts ) {
        var type = record.data.Target_type;
        if (type == 'WIKI') {
            var wstore = this.getStore('Wikis');
            var record = wstore.getById(record.data.Target_id);
            if (record) {
                this.getController('Wikis').openWiki(null, record);
            }
            
        }
        if (type == 'PAGE') {
            // We will have to build treestore and the record, then open the window with them...
            var wstore = this.getStore('Wikis');
            var wiki_id = record.data.Target_id.split('/')[0];
            var parent_id = record.data.Target_id.split('/')[1];
            var page_id = record.data.Target_id.split('/')[2];

            var wikirecord = wstore.getById(wiki_id);

            var treeStore = Ext.create('IW.store.Pages');
            if (parent_id != "") {
                treeStore.setRootNode({
                    expanded: true,
                    Wiki_id: wikirecord.data.Wiki_id,
                    Page_id: parent_id
                });
            } else {
                treeStore.setRootNode({
                    expanded: true,
                    Wiki_id: wikirecord.data.Wiki_id
                });
            }
            
            var me = this;

            treeStore.load({
                scope: this,
                callback: function(records, operation, success) {
                    var pagerecord;

                    var root = treeStore.getRootNode();

                    // treeStore doesn't have proper search method, yet
                    // cascadeBy walks the tree until we get match
                    root.cascadeBy(function (childNode) {
                        if (childNode.get('Page_id') == page_id)
                        {
                            pagerecord = childNode;
                        }
                    }, this);

                    // If we found one, open it!
                    if (pagerecord) {
                        var panel = {
                            store: treeStore
                        };

                        this.getController('Pages').openPage(panel, pagerecord, null, null, null, null);
                    }
                }
            });
        }
        if (type == 'ATTACHMENT') {
            var wiki_id = record.data.Target_id.split('/')[0];
            var att_id = record.data.Target_id.split('/')[1];
            window.open('/att/'+wiki_id+'/'+att_id);
        }
    }
});


