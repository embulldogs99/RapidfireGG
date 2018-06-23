import requests
import json
import time
import psycopg2

def statspull(ep):
    epicusername=ep

    t=requests.Session()

    #api key given to rapidfire.gg
    api = {'TRN-Api-Key':'703cb7b0-4c42-444b-a485-379ed15319b8'}
    # pass api key as header
    r=t.get('http://api.fortnitetracker.com/v1/profile/xbl/'+epicusername, headers = api)
    store=json.loads(r.text)

    squadkill=store['stats']['p9']['kills']['valueInt']
    squadmatch=store['stats']['p9']['matches']['valueInt']
    duokill=store['stats']['p10']['kills']['valueInt']
    duomatch=store['stats']['p10']['matches']['valueInt']
    solokill=store['stats']['p2']['kills']['valueInt']
    solomatch=store['stats']['p2']['matches']['valueInt']
    kills=squadkill+duokill+solokill
    curtime=time.time()
    matchcount=squadmatch+duomatch+solomatch

    print(' ')
    print('--------------------------------------------------')
    print(epicusername+"'s Initial Stats")
    print(' Kills: '+str(kills))
    print(' Matches: '+str(matchcount))
    print('--------------------------------------------------')
    print(' ')

    return (kills,matchcount,curtime)




def playerlist(tournament,rn):
    #########################################################
    ##############  Database Connection   ###################
    conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
    cur = conn.cursor()
    # execute a statement
    cur.execute("SELECT tournaments.epicusername, tournaments.gamertag FROM rfgg.tournaments WHERE tournament='freeweekly2' AND roundnum=1;".format(tournament,rn))
    conn.commit()

    rows = cur.fetchall()
    return rows

    # close the communication with the PostgreSQL
    cur.close()
    conn.close()




conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
cur = conn.cursor()
cur.execute("CREATE TABLE rfgg.tourney_temp (epicusername VARCHAR(500),kills INTEGER,matches INTEGER, time_stamp time);")
conn.commit()

for p in playerlist('freeweekly2',1):
    stats=statspull(p)
    for k,m,c in stats:
        cur.execute("INSERT INTO rfgg.tourney_temp (epicusername,kills,matches,time_stamp) values('{0}','{1}');".format(p,k,m,c))
        conn.commit()

for r in range (1,10):
    cur.execture("SELECT DISTINCT epicusername FROM rfgg.tourney_temp;")
    conn.commit()
    playerlist = cur.fetchall()
    for p in playerlist:
        cur.execture("SELECT epicusername,kills,matches,time_stamp FROM rfgg.tourney_temp WHERE epicusername='{0}';".format(p))
        conn.commit()
        oldstats = cur.fetchall()
        newstats=statspull(p)
        for k,m,c in oldstats:
            for kn,mn,cn in newstats:
                if m<mn:
                    print(p+'has completed a tournament with '+k+' kills')
                    cur.execute("UPDATE rfgg.tournaments (kills,matches) values('{0}','{1}') WHERE tournament='{2}' AND roundnum='{3}' AND gametype='squad';".format(kn,mn,'freeweekly2',1))
                    conn.commit()
                    cur.execute("DELETE FROM rfgg.tourney_temp where epicusername='{0}';".format(e))
                    conn.commit()
                else:
                    time.sleep(10)





cur.execute("DROP TABLE rfgg.tourney_temp;")
conn.commit()
