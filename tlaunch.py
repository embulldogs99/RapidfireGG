import requests
import json
import time
import psycopg2
import cgi
import sys

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
    curtime=int(time.time())
    matchcount=squadmatch+duomatch+solomatch

    print(' ')
    print('--------------------------------------------------')
    print(epicusername+"'s Stats")
    print(' Kills: '+str(kills))
    print(' Matches: '+str(matchcount))
    print('--------------------------------------------------')
    print(' ')

    a={kills,matchcount,curtime}
    return (a)




def playerlist(tournament,rn):
    #########################################################
    ##############  Database Connection   ###################
    teamname=sys.stdin.read()
    teamnamestring=str(teamname)
    print(teamnamestring)

    conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
    cur = conn.cursor()
    # execute a statement
    cur.execute("SELECT tournaments.epicusername, tournaments.gamertag FROM rfgg.tournaments WHERE tournament='{0}' AND roundnum='{1}' AND teamname='{2}';".format(tournament,rn,teamnamestring))
    conn.commit()

    rows = cur.fetchall()


    # close the communication with the PostgreSQL
    cur.close()
    conn.close()
    return rows


######################################################
###### Tourney Run Code ##########################

conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
cur = conn.cursor()

teamname=sys.stdin.read()
teamnamestring=str(teamname)
print(teamnamestring)
sqlstatement="CREATE TABLE rfgg.tourney_temp_{0} (epicusername VARCHAR(500),kills INTEGER,matches INTEGER, time_stamp BIGINT);".format(teamnamestring)
print(sqlstatement)

cur.execute("CREATE TABLE rfgg.tourney_temp_{0} (epicusername VARCHAR(500),kills INTEGER,matches INTEGER, time_stamp BIGINT);".format(teamnamestring))
conn.commit()



for p,rn in playerlist('freeweekly2',1):
    kv,mv,cv =statspull(p)
    cur.execute("INSERT INTO rfgg.tourney_temp_{0} (epicusername,kills,matches,time_stamp) values('{1}','{2}','{3}','{4}');".format(teamnamestring,p,kv,mv,cv))
    conn.commit()
    print('Loading Initial Stats')


for r in range (1,1):
    cur.execute("SELECT epicusername, kills FROM rfgg.tourney_temp_{0};".format(teamnamestring))
    conn.commit()
    playerlist = cur.fetchall()
    for p,t in playerlist:
        if len(p)>3:
            cur.execute("SELECT kills,matches,time_stamp FROM rfgg.tourney_temp_{0} WHERE epicusername='{1}';".format(teamnamestring,p))
            conn.commit()
            rows = cur.fetchall()
            for k,m,c in rows:
                kn,mn,cn = statspull(p)
                if mn<m+3:
                    if round(int(m),0)+.1<round(int(mn),0):
                        print(p+' has submitted tournament entry')
                        cur.execute("UPDATE rfgg.tournaments SET kills='{0}',matches='{1}' WHERE tournament='{2}' AND roundnum='{3}' AND gametype='squad' AND epicusername='{4}';".format((kn-k),(mn-m),'freeweekly2',1,p))
                        conn.commit()
                        cur.execute("DELETE FROM rfgg.tourney_temp_{0} where epicusername='{1}';".format(teamnamestring,p))
                        conn.commit()
                    else:
                        print(p+' is Still at '+str(mn)+' Matches')
                        time.sleep(2)
                else:
                    time.sleep(2)
        else:
            time.sleep(10)
    time.sleep(60)


cur.execute("DROP TABLE rfgg.tourney_temp_{0};".format(teamnamestring))
conn.commit()
cur.close()
conn.close()
