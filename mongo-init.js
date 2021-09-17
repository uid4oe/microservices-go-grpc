db = db.getSiblingDB('oe');
db.createUser({
  user: 'oe',
  pwd: 'oe',
  roles: [{
    role: 'readWrite',
    db: 'oe'
  }]
})
