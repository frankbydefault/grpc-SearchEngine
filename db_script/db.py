from pydal import DAL, Field

db = DAL('mysql://searchengine:S34rch3r_3ng1n3@localhost/search_engine', folder="databases")

db.define_table('data',
    Field('title', 'string'),
    Field('description', 'string')
)

db.define_table('keywords',
    Field('id_data', 'integer'),
    Field('keyword', 'string')
)

db.commit()