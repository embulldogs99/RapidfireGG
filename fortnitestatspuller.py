import requests
import json
import time
import psycopg2


def statspull(e,ep,c):
    email=e
    epicusername=ep
    console=c
    t=requests.Session()

    #api key given to rapidfire.gg
    api = {'TRN-Api-Key':'703cb7b0-4c42-444b-a485-379ed15319b8'}
    # pass api key as header

    r=t.get('http://api.fortnitetracker.com/v1/profile/'+console+'/'+epicusername, headers = api)
    store=json.loads(r.text)


    squadkill=store['stats']['curr_p9']['kills']['valueInt']
    squadmatch=store['stats']['curr_p9']['matches']['valueInt']
    squadkm=round(squadkill/float(squadmatch),2)
    duokill=store['stats']['curr_p10']['kills']['valueInt']
    duomatch=store['stats']['curr_p10']['matches']['valueInt']
    duokm=round(duokill/float(duomatch),2)
    solokill=store['stats']['p2']['kills']['valueInt']
    solomatch=store['stats']['p2']['matches']['valueInt']
    solokm=round(solokill/float(solomatch),2)

    lastupdated=time.time()

    #########################################################
    ##############  Database Connection   ###################
    conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
    cur = conn.cursor()
    # execute a statement
    cur.execute("INSERT INTO rfgg.fortniteplayerstats (last_updated,console,squadkill,squadmatch,squadkm,duokill,duomatch,duokm,solokill,solomatch,solokm, email, epicusername) VALUES (%s, %s, %s, %s, %s,%s,%s,%s,%s,%s,%s,%s,%s)", (lastupdated,console,squadkill,squadmatch,squadkm,duokill,duomatch,duokm,solokill,solomatch,solokm,email,epicusername))
    conn.commit()
    # close the communication with the PostgreSQL
    cur.close()
    conn.close()



def xblpull():
    #########################################################
    ##############  Database Connection   ###################
    conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
    cur = conn.cursor()
    # execute a statement
    cur.execute("SELECT DISTINCT members.email, members.epicusername FROM rfgg.members WHERE members.epicusername<>'?';")
    conn.commit()

    rows = cur.fetchall()
    for x,y in rows:
        statspull(x,y,'xbl')
    # close the communication with the PostgreSQL
    cur.close()
    conn.close()


xblpull()
