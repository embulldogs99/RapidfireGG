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
    wins=0

    print(' ')
    print('--------------------------------------------------')
    print(epicusername+"'s Initial Stats")
    print(' Squad Kills: '+str(squadkill))
    print(' Squad Matches: '+str(squadmatch))
    print(' Duo Kills: '+str(duokill))
    print(' Duo Matches: '+str(duomatch))
    print(' Solo Kills: '+str(solokill))
    print(' Solo Matches: '+str(solomatch))
    print('--------------------------------------------------')
    print(' ')


    timestart=time.time()
    matchcount=0

    return (kills,matchcount,timestart)


    #
    # for x in range(1,60):
    #
    #     s=requests.Session()
    #
    #     # pass api key as header
    #     r=s.get('http://api.fortnitetracker.com/v1/profile/xbl/'+epicusername, headers = api)
    #     newstore=json.loads(r.text)
    #
    #     print('--------------------------------------------------')
    #     print(str(round((time.time()-timestart)/60,0))+' Mins after Tournament Start')
    #     print('Squad Kills: '+str(round(newstore['stats']['p9']['kills']['valueInt']-squadkillstart)))
    #     print('Squad Matches: '+str(round(newstore['stats']['p9']['matches']['valueInt']-squadmatchstart)))
    #     print('Duo Kills: '+str(round(newstore['stats']['p10']['kills']['valueInt']-duokillstart)))
    #     print('Duo Matches: '+str(round(newstore['stats']['p10']['matches']['valueInt']-duomatchstart)))
    #     print('Solo Kills: '+str(round(newstore['stats']['p2']['kills']['valueInt']-solokillstart)))
    #     print('Solo Matches: '+str(round(newstore['stats']['p2']['matches']['valueInt']-solomatchstart)))
    #     print('--------------------------------------------------')
    #     print(' ')
    #
    #     kills=round(newstore['stats']['p9']['kills']['valueInt']-squadkillstart)+round(newstore['stats']['p10']['kills']['valueInt']-duokillstart)+round(newstore['stats']['p2']['kills']['valueInt']-solokillstart)
    #     wins=0
    #
    #
    #     curtime=str(round((time.time()-timestart)/60,0))+' Minutes'
    #     squadmatchcount=round(newstore['stats']['p9']['matches']['valueInt']-squadmatchstart)
    #     duomatchcount=round(newstore['stats']['p10']['matches']['valueInt']-duomatchstart)
    #     solomatchcount=round(newstore['stats']['p2']['matches']['valueInt']-solomatchstart)
    #     matchcount=squadmatchcount+duomatchcount+solomatchcount
    #
    #     if matchcount is 1:
    #         #########################################################
    #         ##############  Database Connection   ###################
    #         conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
    #         cur = conn.cursor()
    #         # execute a statement
    #         cur.execute("INSERT INTO rfgg.tournaments WHERE tournament=%s AND roundnum=%s AND gametype=%s AND epicusername=%s (tournament,roundnum,gametype,epicusername,wins,kills,matches,teamname) VALUES (%s, %s, %s, %s, %s,%s,%s,%s)", (tournament,roundnum,gametype,tournament,roundnum,gametype,gamertag,epicusername,wins,kills,matchcount))
    #         print(epicusername+" Finished Round "+roundnum)
    #         conn.commit()
    #         # close the communication with the PostgreSQL
    #         cur.close()
    #         conn.close()
    #
    #     else:
    #         time.sleep(10)


def playerpull(tournament,rn):
    #########################################################
    ##############  Database Connection   ###################
    conn = psycopg2.connect("dbname='postgres' user='postgres' password='rk' host='localhost' port='5432'")
    cur = conn.cursor()
    # execute a statement
    cur.execute("SELECT tournaments.epicusername, tournaments.gamertag FROM rfgg.tournaments WHERE tournament='freeweekly2' AND roundnum='{1}';".format(tournament,rn))
    conn.commit()

    rows = cur.fetchall()
    for x,y in rows:
        return x

    # close the communication with the PostgreSQL
    cur.close()
    conn.close()



original=statspull(playerpull('freeweekly2',1))
print(original)
