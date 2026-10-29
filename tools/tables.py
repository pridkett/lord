import re, sys
src = open(sys.argv[1], encoding='latin-1').read().split('\n')
def block(start_pat):
    i = next(n for n,l in enumerate(src) if l.startswith(start_pat))
    out=[]; j=i+1
    while not src[j].startswith('];'):
        out.append(src[j]); j+=1
    return '\n'.join(out)
fieldmap = {'name':'Name','str':'Str','gold':'Gold','weapon':'Weapon','exp':'Exp','hp':'HP','death':'Death',
 'is_arena':'IsArena','swear':'Swear','needstr1':'NeedStr1','needstr2':'NeedStr2','hp_gained':'HPGained',
 'str_gained':'StrGained','def':'Def','need':'Need','price':'Price','num':'Num'}
def jsstr(v):
    v = v.strip()
    q=v[0]; body=v[1:-1]
    body = body.replace("\\'", "'")
    return '"' + body.replace('\\','\\\\').replace('"','\\"') + '"' if '\\' not in v[1:-1].replace("\\'","") else '"' + body.replace('"','\\"') + '"'
def objs(text):
    res=[]
    for m in re.finditer(r'\{(.*?)\}', text, re.S):
        body=m.group(1); fields=[]
        for fm in re.finditer(r"([a-z_0-9]+)\s*:\s*('(?:[^'\\]|\\.)*'|[^,\n}]+)", body):
            k,v=fm.group(1),fm.group(2).strip()
            if v.startswith("'"): v=jsstr(v)
            fields.append('%s: %s'%(fieldmap[k],v))
        res.append(fields)
    return res
def emit(name, typ, text):
    print('var %s = []%s{'%(name,typ))
    for f in objs(text):
        print('\t{%s},'%', '.join(f))
    print('}\n')
print('// Code generated from ref/synchronet/lord.js by tools/tables.py; DO NOT EDIT.\n\npackage game\n')
emit('trainerStats','Player',block('var trainer_stats'))
emit('monsterStats','Player',block('var monster_stats'))
emit('armourStats','WeaponArmour',block('var armour_stats'))
emit('weaponStats','WeaponArmour',block('var weapon_stats'))
