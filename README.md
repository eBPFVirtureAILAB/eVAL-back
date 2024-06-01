# eVAL-backend server

## 

## 개발 테스트 환경 구축 방식

### 배포

backend (go) - database (postgresql)

* requirement

  * docker

~~~ shell

# 환경 설정 .env 파일로 커스텀
# global 폴더에서 커스텀 항목 찾을 수 있음
$ vim .env
...

# 환경 배포
$ docker compose -f dockercompose.yaml up -d

# 배포 확인
$ docker ps

# 백엔드 서버 로그 확인
$ docker logs -f eval-backend-main-1

# postgresql db 확인
$ docker exec -it eval-backend-db-1 /bin/bash
-> $ psql -U postgres -d eval
...

~~~

### 테스트를 위한 쿠버네티스 환경 필요
DB 초기화 과정 필요
- user 등록
- kubeconfig DB에 등록
- ...
