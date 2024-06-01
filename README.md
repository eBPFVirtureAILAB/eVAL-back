# eVAL-backend server

## 

## 개발 테스트 환경 구축 방식

### 배포
backend image는 dockerfile로 빌드

이후 docker compose 를 사용
=> backend (go) - database (postgresql)

### 테스트를 위한 쿠버네티스 환경 필요
DB 초기화 과정 필요
- user 등록
- kubeconfig DB에 등록
- ...
