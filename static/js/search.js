var index = elasticlunr(function () {
  this.addField('name');
  this.addField('categories');
  this.addField('description');
  this.addField('services');
  this.addField('email');
  this.addField('address');
  this.addField('languages');
  this.addField('population_served');
  this.addField('website');
  this.addField('business_line');
  this.addField('crisis_line');
  this.addField('fax');
});

for (var i = 0; i < resources.length; i++) {
    resources[i].id = i;
    index.addDoc(resources[i]);
}
